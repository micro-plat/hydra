###  1. 注册中心[hydra_registry_info]

| 字段名       | 类型           | 默认值  | 为空  |  约束  | 描述       |
| ------------ | -------------- | :-----: | :---: | :----: | :--------- |
| id           | id             |   100   |  否   | PK,seq | 编号       |
| path         | varchar2(64)   |         |  否   |   UNQ     | 路径       |
| temp         | number(1)      |    0    |  否   |        | 临时节点   |
| value        | varchar2(1024) |         |  否   |        | 内容       |
| data_version | number(20)     |         |  是   |        | 数据版本号 |
| acl_version  | number(20)     |         |  是   |        | 访问版本号 |
| create_time  | date           | sysdate |  否   |   d    | 创建时间   |
| update_time  | date           | sysdate |  否   |   d    | 更新时间   |