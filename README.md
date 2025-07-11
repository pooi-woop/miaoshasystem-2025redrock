（readme是给AI写的）

# 红岩网校 2025 年春季考核作业：秒杀解决方案

## 项目简介

这是一个秒杀系统的解决方案，旨在应对高并发场景下的秒杀业务。通过合理的技术选型和架构设计，解决了高并发处理、库存超卖以及订单存储等问题，同时提供了完整的接口文档供开发者使用。

## 项目特点

1. **高并发处理**：通过中间件限制高频请求，结合 Kafka 消息队列进行缓冲，有效应对高并发请求。
2. **库存管理**：优先进行减少库存操作，防止超卖，确保业务稳定性。
3. **订单存储**：使用 Redis 存储键值对，提升订单存储效率，简化操作流程。
4. **完整的接口文档**：提供了详细的接口说明，方便开发者快速上手和集成。

## 技术栈

- **web框架**：gin
- **消息队列**：Kafka
- **数据库**：Redis，mysql
- **身份认证**：JWT

## 项目结构

```plaintext
miaoshasystem-2025redrock/
├── src/
│   ├── controllers/       # 控制器层代码
│   ├── models/            # 数据模型定义
│   ├── services/          # 业务逻辑层代码
│   ├── middlewares/       # 中间件代码
│   ├── routes/            # 路由定义
│   └── utils/             # 工具类代码
├── config/                # 配置文件
├── tests/                 # 测试代码
├── package.json           # 项目依赖配置
├── README.md              # 项目说明文档
└── .env                   # 环境变量配置文件
```

## 安装与运行

### 环境准备

1. **Node.js**：请确保已安装 Node.js 环境。
2. **Kafka**：安装并启动 Kafka 服务，确保其正常运行。
3. **Redis**：安装并启动 Redis 服务，用于订单存储。

### 安装依赖

```bash
npm install
```

### 配置环境变量

复制 `.env.example` 文件并重命名为 `.env`，根据实际情况修改其中的配置项，例如 Kafka 和 Redis 的连接信息。

### 启动项目

```bash
npm start
```

项目将启动并监听默认端口 [端口号]，你可以通过访问相关接口进行测试。

## 接口文档

### 1. 注册接口

- **接口地址**：`POST /register`
- **请求参数**：
  | 参数名   | 类型   | 是否必填 | 示例值                     | 描述       |
  |----------|--------|----------|----------------------------|------------|
  | name     | string | 是       | "张三"                     | 用户名     |
  | age      | int    | 是       | 20                         | 年龄       |
  | address  | string | 是       | "北京市海淀区"             | 地址       |
  | avatar   | string | 是       | "http://example.com/avatar.jpg" | 头像链接   |
  | password | string | 是       | "123456"                   | 密码       |
- **请求体示例**：
  ```json
  {
      "name": "张三",
      "age": 20,
      "address": "北京市海淀区",
      "avatar": "http://example.com/avatar.jpg",
      "password": "123456"
  }
  ```
- **返回结果**：
  | 状态码 | 返回值                                 | 描述               |
  |--------|----------------------------------------|--------------------|
  | 200    | {"message": "user registered successfully"} | 注册成功           |
  | 400    | {"error": "error message"}             | 请求参数错误       |
  | 500    | {"error": "user creation failed"}      | 服务器内部错误     |

### 2. 登录接口

- **接口地址**：`POST /login`
- **请求参数**：
  | 参数名   | 类型   | 是否必填 | 示例值   | 描述   |
  |----------|--------|----------|----------|--------|
  | name     | string | 是       | "张三"   | 用户名 |
  | pass     | string | 是       | "123456" | 密码   |
- **请求体示例**：
  ```json
  {
      "name": "张三",
      "pass": "123456"
  }
  ```
- **返回结果**：
  | 状态码 | 返回值                     | 描述               |
  |--------|----------------------------|--------------------|
  | 200    | {"token": "token_string"}  | 登录成功，返回JWT令牌 |
  | 400    | {"error": "error message"} | 请求参数错误       |
  | 401    | {"error": "Invalid token"} | 令牌无效           |

### 3. 创建秒杀产品接口

- **接口地址**：`POST /createmiaosha`
- **请求参数**：
  | 参数名             | 类型   | 是否必填 | 示例值                     | 描述               |
  |--------------------|--------|----------|----------------------------|--------------------|
  | name               | string | 是       | "iPhone 13"                | 产品名称           |
  | num                | int    | 是       | 100                        | 产品数量           |
  | producter          | string | 是       | "Apple"                    | 生产商             |
  | time_begintokill   | int64  | 是       | 1640995200                 | 秒杀开始时间（时间戳） |
  | time_endkill       | int64  | 是       | 1641081600                 | 秒杀结束时间（时间戳） |
- **请求体示例**：
  ```json
  {
      "name": "iPhone 13",
      "num": 100,
      "producter": "Apple",
      "time_begintokill": 1640995200,
      "time_endkill": 1641081600
  }
  ```
- **返回结果**：
  | 状态码 | 返回值                                 | 描述               |
  |--------|----------------------------------------|--------------------|
  | 200    | {"message": "Product creation failed"} | 创建成功           |
  | 400    | {"error": "error message"}             | 请求参数错误       |
  | 500    | {"error": "Product creation failed"}   | 服务器内部错误     |

### 4. 秒杀接口

- **接口地址**：`POST /miaosha/:productName`
- **请求参数**：
  | 参数名   | 类型   | 是否必填 | 示例值                     | 描述       |
  |----------|--------|----------|----------------------------|------------|
  | productName | string | 是       | "iPhone 13"                | 产品名称   |
  | token    | string | 是       | "token_string"              | JWT令牌   |
- **请求体示例**：
  ```json
  {
      "token": "token_string"
  }
  ```
- **返回结果**：
  | 状态码 | 返回值                                                                 | 描述               |
  |--------|------------------------------------------------------------------------|--------------------|
  | 200    | {"success": "订单创建成功", "time": "current_time", "username": "username", "product name": "productName", "注意": "未支付的订单将在一个小时之后失效"} | 秒杀成功           |
  | 400    | {"error": "error message"}                                             | 请求参数错误       |
  | 401    | {"error": "Invalid token"}                                             | 令牌无效           |
  | 500    | {"error": "Failed to send request to Kafka"}                           | 服务器内部错误     |

### 5. Kafka 消费者接口（内部接口，不对外暴露）

- **功能描述**：该接口用于处理 Kafka 队列中的秒杀请求，将请求发送到后端进行处理。
- **请求参数**：无
- **返回结果**：无

## 测试

项目中包含测试代码，你可以通过以下命令运行测试：

```bash
npm test
```

## 贡献指南

欢迎贡献代码！如果你有任何改进意见或修复了问题，请按照以下步骤操作：

1. **Fork** 本项目到你的 GitHub 账号。
2. 创建一个新的分支：`git checkout -b feature/your-feature-name`
3. 进行代码修改并提交：`git commit -m "Add your feature"`
4. 将代码推送到你的分支：`git push
