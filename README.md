# forum-be
木犀论坛后端仓库

mysql

websocket: 聊天

拆分成微服务microservice，使用grpc（protobuf）通信：
- gateway
- post
- user
- chat
- feed

casbin: 鉴权

redis: 点赞; 聊天; tag缓存; tag排行榜; 消息队列(pub-sub); hot post

log: zap

CI: .drone.yml

Kafka: 动态功能的消息队列
