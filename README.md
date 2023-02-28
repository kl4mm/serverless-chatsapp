# serverless-chatsapp
Chat app built with go/lambda and deployed using terraform, initially following 
[this](https://docs.aws.amazon.com/apigateway/latest/developerguide/websocket-api-chat-app.html#websocket-api-chat-app-create-dependencies) guide,
but will extend in future to use rooms like [chatsapp](https://github.com/7junky/chatsapp).

## Issues

* Need to redeploy the API after it provisions for the first time or the connect lambda doesn't run at all
* Posting to connections results in this error: `Post "https://execute-api.eu-west-2.amazonaws.com/@connections/BEGNudC9rPECF2w%3D": dial tcp: lookup execute-api.eu-west-2.amazonaws.com on 169.254.78.1:53: no such host`
