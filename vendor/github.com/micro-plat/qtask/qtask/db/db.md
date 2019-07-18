## 数据表

### 任务表[tsk_system_task]

| 字段名            | 类型          | 默认值  | 为空 |  约束  | 描述                               |
| ----------------- | ------------- | :-----: | :--: | :----: | :--------------------------------- |
| task_id           | number(20)    |         |  否  | PK,SEQ | 编号                               |
| name              | varchar2(32)  |         |  否  |        | 名称                               |
| create_time       | date          | sysdate |  否  |        | 创建时间                           |
| last_execute_time | date          |         |  是  |        | 上次执行时间                       |
| next_execute_time | date          |         |  否  |        | 下次执行时间                       |
| max_execute_time  | date          |         |  否  |        | 执行期限(此时间前的任务可以被执行) |
| next_interval     | number(10)    |         |  否  |        | 时间间隔,秒数                      |
| count             | number(10)    |    0    |  否  |        | 执行次数                           |
| status            | number(2)     |         |  否  |        | 状态(20 等待，30 正在,0 已处理)    |
| batch_id          | number(20)    |         |  是  |        | 执行批次号                         |
| queue_name        | varchar2(64)  |         |  否  |        | 消息队列                           |
| msg_content       | varchar2(256) |         |  是  |        | 消息内容                           |

### MYSQl 序列表[tsk_system_seq]

| 字段名      | 类型         | 默认值  | 为空 |  约束  | 描述     |
| ----------- | ------------ | :-----: | :--: | :----: | :------- |
| seq_id      | number(20)   |         |  否  | PK,SEQ | 编号     |
| name        | varchar2(32) |         |  否  |        | 名称     |
| create_time | date         | sysdate |  否  |        | 创建时间 |
