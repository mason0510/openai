package gpt

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"openai/config"
	"strings"
	"sync"
	"time"
)

const (
	api = "https://api.openai.com/v1/completions"
)

var (
	// 结果缓存（主要用于超时，用户重新提问后能给出答案）
	resultCache sync.Map
)

type response struct {
	ID string `json:"id"`
	// Object  string                 `json:"object"`
	// Created int                    `json:"created"`
	// Model   string                 `json:"model"`
	Choices []choiceItem `json:"choices"`
	// Usage   map[string]interface{} `json:"usage"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type choiceItem struct {
	Text string `json:"text"`
	// Index        int    `json:"index"`
	// Logprobs     int    `json:"logprobs"`
	// FinishReason string `json:"finish_reason"`
}

// OpenAI可能无法在希望的时间内做出回复
// 使用goroutine + channel 的形式，不管是否能及时回复用户，后台都打印结果
func Query(isFast bool, msg string, timeout time.Duration) string {
	start := time.Now()
	ch := make(chan string, 1)
	ctx, candel := context.WithTimeout(context.Background(), timeout)
	defer candel()

	cacheVal, ok := resultCache.Load(msg)
	if ok {
		return cacheVal.(string)
	}

	go func() {
		defer close(ch)
		result, err := completions(isFast, msg, time.Second*100)
		if err != nil {
			result = "发生错误「" + err.Error() + "」，您重试一下"
		}
		ch <- result
		// 超时，内容未通过接口及时回复，打印内容及总用时
		since := time.Since(start)
		if since > timeout {
			// TODO定时清理
			resultCache.Store(msg, result)
			log.Printf("超时%ds，「%s」-「%s」\n", int(since.Seconds()), msg, result)
		}
	}()

	var result string
	select {
	case result = <-ch:
	case <-ctx.Done():
		result = "超时啦，你等下再问我一遍，我一定告诉你！"
	}

	log.Printf("用时%ds，「%s」-「%s」\n", int(time.Since(start).Seconds()), msg, result)

	return result
}

func getFromCache(msg string) string {
	v, ok := resultCache.Load(msg)
	if ok {
		return v.(string)
	}
	return ""
}

// https://beta.openai.com/docs/api-reference/making-requests
func completions(isFast bool, msg string, timeout time.Duration) (string, error) {
	wordSize := 1024 // 中文字符数量
	temperature := 0.5

	if !isFast {
		wordSize = 4000
		temperature = 0.95
	}

	// start := time.Now()
	params := map[string]interface{}{
		"model":  "text-davinci-003",
		"prompt": msg,
		// 影响回复速度和内容长度。小则快，但内容短，可能是截断的。
		"max_tokens": wordSize,
		// 0-1，默认1，越高越有创意
		"temperature":       temperature,
		"top_p":             1,
		"frequency_penalty": 0,
		"presence_penalty":  0,
		"stop":              "。",
	}

	bs, _ := json.Marshal(params)

	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequest("POST", api, bytes.NewReader(bs))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.ApiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data response
	json.Unmarshal(body, &data)
	if len(data.Choices) > 0 {
		result := strings.TrimPrefix(data.Choices[0].Text, "？")
		return strings.TrimSpace(result), nil
	}

	return data.Error.Message, nil
}
