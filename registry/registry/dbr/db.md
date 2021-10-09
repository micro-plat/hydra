<!--
 * @Description: 
 * @Autor: taoshouyin
 * @Date: 2021-09-18 09:36:32
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-18 09:38:51
-->
###  1. 注册中心[hydra_registry_info]

| 字段名       | 类型           | 默认值  | 为空  |  约束  | 描述       |
| ------------ | -------------- | :-----: | :---: | :----: | :--------- |
| id           | number(10)     |   100   |  否   | PK,seq | 编号       |
| path         | varchar2(256)   |         |  否   |  UNQ   | 路径       |
| value        | varchar2(4096) |         |  否   |        | 内容       |
| is_temp      | number(1)      |    0    |  否   |        | 临时节点   |
| is_delete    | number(1)      |    1    |  否   |        | 已删除     |
| data_version | number(20)     |         |  是   |        | 数据版本号 |
| create_time  | date           | sysdate |  否   |   d    | 创建时间   |
| update_time  | date           | sysdate |  否   |   d    | 更新时间   |