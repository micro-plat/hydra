#### 1. 用户信息表(ucs_user_info)

| 字段名          | 类型（长度）  | 可为空 | 默认    | 约束   | 说明                                                |
| --------------- | ------------- | ------ | ------- | ------ | --------------------------------------------------- |
| user_id         | number(20)    | 否     |         | PK     | 编号                                                |
| app_id          | varchar2(32)  | 否     |         | IS,SEQ | 应用appid                                           |
| open_id         | varchar2(32)  | 否     |         | IS     | 用户openid                                          |
| subscribe       | number(1)     | 否     |         | IS     | 关注状态(0-关注,1-未关注)                           |
| status          | number(1)     | 否     | 0       | IS     | 用户状态（0-启用，9-禁用，1-锁定）                  |
| create_time     | date          | 否     | sysdate | IS     | 创建时间                                            |
| nick_name       | varchar2(32)  | 是     |         | IS     | 用户昵称                                            |
| gender          | number(1)     | 是     |         | IS     | 性别(0-女,1-男,3-未知)                              |
| subscribe_time  | date          | 是     |         | IS     | 关注时间                                            |
| update_time     | date          | 是     |         | IS     | 修改时间                                            |


#### 2. app基本信息表（ucs_app_info）

| 字段名            | 类型（长度）   | 可为空 | 默认    | 约束 | 说明                                               |
| ----------------- | -------------- | ------ | ------- | ---- | -------------------------------------------------- |
| app_id            | varchar2(32)   | 否     |         | PK   | 应用appid                                          |
| app_raw_id        | varchar2(32)   | 是     |         | IS   | 应用原始ID                                         |
| app_name          | varchar2(32)   | 否     |         | IS   | 应用名称                                           |
| plat_type         | varchar2(16)   | 否     |         | IS   | 平台类型(wechat-微信,alipay-支付宝,bestpay-翼支付) |
| app_type          | number(1)      | 否     |         | IS   | 应用类型(1-公众号,2-小程序)                        |
| app_status        | number(1)      | 否     | 1       | IS   | 应用状态(1-启用,0-禁用)                            |
| create_time       | date           | 否     | sysdate | IS   | 创建时间                                           |
