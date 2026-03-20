# 项目名称（Boilerplate / Lottery Service）

一个基于 Go 语言开发的后端服务，采用分层架构设计，集成 Gin、MySQL(GORM)、Redis、Viper，具备基础的业务开发能力与工程规范。

---

## 🚀 技术栈

* **Web框架**：Gin
* **ORM**：GORM
* **数据库**：MySQL
* **缓存**：Redis（go-redis）
* **配置管理**：Viper
* **认证方式**：JWT
* **日志**：标准库 log（可扩展）

---

## 📂 项目结构

```
├── cmd/                # 程序入口
├── config/             # 配置加载（Viper）
├── internal/
│   ├── dto/            # 请求参数绑定结构
│   ├── handler/        # HTTP处理层（Controller）
│   ├── service/        # 业务逻辑层
│   ├── repository/     # 数据访问层（DAO）
│   ├── model/          # 数据库实体（GORM模型）
│   ├── router/         # 路由注册
│   └── middleware/     # 中间件（JWT、日志等）
├── pkg/
│   ├── constants/      # 常量
│   ├── jwt/            # JWT工具封装
│   ├── redis/          # Redis工具封装
│   ├── mysql/          # 数据库初始化
│   └── response/       # 统一响应结构
└── go.mod
```

---

## 🧱 分层架构说明

本项目采用经典的三层架构：

### 1️⃣ Handler 层（接口层）

* 负责 HTTP 请求处理
* 参数绑定（Gin binding）
* 调用 Service 层
* 返回统一 JSON 响应

```go
func (l *LotteryHandler) Add(c *gin.Context) {
    var form dto.LotteryAddReq
    if err := c.ShouldBindJSON(&form); err != nil {
        c.JSON(http.StatusBadRequest, response.FailWithMsg("ruleId is required"))
        return
    }
    id := c.MustGet("userId")
    res, err := redis.Lock(c, constants.LotteryUserLock, func() (r *response.CommonResult) {
        reward, err := service.NewLotteryService().Add(c, &form)
        if err != nil {
            log.Printf("抽奖失败，用户ID: %v，错误: %v", id, err)
            c.JSON(http.StatusInternalServerError, response.Fail())
            return
        }
        return response.Success(reward)
    })
    if err != nil {
        log.Printf("抽奖失败，用户ID: %v，错误: %v", id, err)
        c.JSON(http.StatusInternalServerError, response.Fail())
        return
    }
    c.JSON(http.StatusOK, res)
}
```

---

### 2️⃣ Service 层（业务层）

* 核心业务逻辑处理
* 事务控制
* 调用 Repository
* 不关心 HTTP 细节

---

### 3️⃣ Repository 层（数据访问层）

* 负责数据库操作（GORM）
* 封装 CRUD
* 与具体表结构强绑定

---

## 🧠 设计要点

### ✅ 1. 配置管理（Viper）

* 支持多环境（dev / prod）
* 统一配置加载

---

### ✅ 2. 数据库模型设计

* 使用 GORM 映射数据库表
* 自定义类型实现“类枚举”（如状态字段）
* 通过指针控制字段是否参与 SQL（实现类似 MyBatis 的动态 SQL）

---

### ✅ 3. Redis 封装

* 封装基础操作
* 支持分布式锁（基于 SET NX + TTL）
* 提供统一调用入口

---

### ✅ 4. 统一响应结构

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

---

### ✅ 5. 中间件

* JWT 鉴权
* 请求拦截
* CORS
* RequestID请求标识

---

## 🔐 分布式锁实现（Redis）

基于 Redis 实现简单分布式锁：

* 使用 `SET key value NX EX`
* 使用 UUID 作为锁标识
* 释放锁时校验 value，避免误删

---

## ⚙️ 配置示例

```yaml
server:
  port: 8080

mysql:
  host: localhost
  port: 3306
  username: username
  password: password
  Database: db

redis:
  host: localhost #地址
  port: 6379 #端口
  password: password
  timeout: 10000 # 连接超时时间（毫秒）
  database: 3
  max-active: 200 # 连接池最大连接数（使用负值表示没有限制）
  max-idle: 10 # 连接池中的最大空闲连接
  min-idle: 0 # 连接池中的最小空闲连接

token:
  # 令牌自定义标识
  header: Authorization
  # 令牌密钥
  secret: secret
  # 令牌有效期
  expireTime: 1440
```

---

## ▶️ 启动项目

```bash
go mod tidy
go run cmd/main.go
```

---

## 📌 后续优化方向

* [ ] 日志系统（zap / logrus）
* [ ] 配置热更新
* [ ] 接口限流（Redis / Token Bucket）
* [ ] 链路追踪（OpenTelemetry）
* [ ] 单元测试完善

---

## ✨ 项目特点

* 清晰的分层结构（Handler / Service / Repository）
* 类似 MyBatis 的动态 SQL 控制（通过指针实现）
* Redis + MySQL 常见后端组合实践
* 适合作为 Go Web 项目基础模板

---

## 📄 License

MIT
