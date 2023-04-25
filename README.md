# lynn-toolkit
一个用于构建golang应用程序工具箱
## 1. 介绍
此工具箱包括了一些常用的golang应用程序开发的工具，如：wire、日志、配置、数据库、缓存、消息队列、定时任务、http服务等。
* 每一个组件都提供了wire的provider，可以直接使用wire进行注入。
* 每一个组件都提供了配置(如果有)Model,可以直接在项目的配置文件中配置，但是得用yaml格式。

## 2. 使用说明
### 2.1 数据库，ORM
#### 2.1.1 gorm
* 在项目的配置文件中配置数据库连接信息
```yaml
# 数据库配置
db:
  dsn: root:aa123123@tcp(127.0.0.1:3306)/liaoma-payment-ylb?charset=utf8mb4&parseTime=True&loc=Local
  # 最大连接数
  maxOpenConn: 5
  # 最大闲置连接数
  maxIdleConn: 1
  ignoreRecordNotFoundError: false
  # 日志等级:1-Silent,2-Error,3-Warn,4-Info
  logLevel: 3
```
* 在项目中使用,
```go
database.NewDBClientWithProfile
```