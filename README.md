# 秒杀系统

## 一、高并发处理
我通过中间件 **ban 掉高频发送请求的脚本哥**，通过**Kafka 消息队列** 的缓冲机制，处理高并发请求


## 二、库存问题
优先进行 **减少库存操作**，这样就能防止超卖。即使后续服务未能成功，也只会少卖不会多卖，从而避免因超卖引发的业务问题。

## 三、订单存储
之前寒假考核的时候，我还不太会用非关系型数据库，存订单相当恼火。这次通过使用 **Redis** 存储键值对，解决了这个问题，爽。
