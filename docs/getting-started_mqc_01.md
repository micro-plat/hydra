# 构建消息消费服务（MQC）

MQC订阅消息队列的topic,收到消息后调用本地服务执行，支持的消息队列有:

|   名称   | 说明              |
| :------: | ----------------- |
| activeMQ | 基于stomp协议实现 |
| rabbitMQ | 基于stomp协议实现 |
|  redis   | 基于list实现      |
|   mqtt   | 物联网消息队列    |


