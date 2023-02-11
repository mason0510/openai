### 最近持续优化中，喜欢的同学给个🌟关注一下
### 介绍
- 功能。通过调用`OpenAI`提供微信公众号自动回复服务。内容可能及其离谱，仅供参考。
- 注意。有别于网页版`ChatGPT`基于GPT-3.5，本项目是调用GPT-3，没那么强大。
- 缺陷。不支持上下文；[微信限制，最久15s做出回复](https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Passive_user_reply_message.html)，可能超时，可以过一会儿再问一下 (缓存结果)  
- 代码。因为功能比较简单，未使用框架，后续可能会优化。如果你有问题或者优化点欢迎联系我探讨，github可能不及时查看，加QQ:`772532526`

### 准备
- [OpenAI账号](https://beta.openai.com) (需要梯子,接口不是完全免费，只是初始赠送金额，具体看官方文档)
- [微信公众号](https://mp.weixin.qq.com/)
- 服务器、域名、`Docker`或`Golang`环境

### 操作流程
1. 获取`API_KEY`。[OpenAI](https://beta.openai.com/account/api-keys) （如果访问被拒绝，注意全局代理，打开调试，Application清除LocalStorage后刷新，实测可以）
2. 获取微信公众号令牌`TOKEN`：[微信公众平台](https://mp.weixin.qq.com/) -> 基本配置 -> 生成令牌 
3. 使用以上两个参数按照下边**项目部署**。
4. 设置公众号服务器地址(端口必须80或者443)，通过nginx代理到`http://127.0.0.1:端口`. 配置举例
  ```conf
  server {
    listen 80;
    server_name xxx.com;

    location / {
      # 略
    }

    # 举例9001端口，公众号服务器地址设置为 http://xxx.com/chatgpt/wx; 
    location /chatgpt/ {
      proxy_pass http://127.0.0.1:9001/;
    }
  }
  ```


### 项目部署
> 举例部署端口为 9001
> 可以不设置 WX_TOKEN，但 API_KEY 是必须的

- Docker
  ```bash
  docker run -p 9001:8080 -e API_KEY=xxx -e WX_TOKEN=xxx -d -v $PWD/log:/app/log tomatocuke/openai
  ```
- Golang运行
  ```bash 
  git clone https://github.com/tomatocuke/chatgpt.git
  cd chatgpt
  go run main.go -PORT=9001 -API_KEY=xxx -WX_TOKEN=xxx 
  ```

### 说明
- 日志。查看`tail -f ./log/chatgpt.log`
- 简单测试。`curl 'http://127.0.0.1:9001/test?msg=中国在哪个洲'` 
- 模式。分两种模式，默认快速回答，但是结果可能不全面。 普通模式，公众号问题中加入`详细`或者`代码`最大化回复，test接口增加`mode=full`参数以调节，例 `curl 'http://127.0.0.1:9001/test?msg=评价一下三体&mode=full'`

这里
export API_KEY=sk-eK9a0n
export WX_TOKEN=123***

