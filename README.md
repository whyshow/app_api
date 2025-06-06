# 用户、新闻资讯 API服务

基于Golang Gin框架构建的RESTful API服务，提供用户认证、新闻数据管理等功能。

## 功能特性

- JWT用户认证体系
- 新闻数据定时抓取与存储
- 分页查询接口
- 密码安全哈希存储
- 配置热更新支持
- 请求日志追踪

## 快速开始

- 首次启动前，请确保数据库已创建并配置正确。
- 启动后将自动创建数据库表
- ```bash   
  # 启动服务
   go run main.go

### 环境要求

- Go 1.23+
- MySQL 5.7+

### 编辑 config.yaml

````
server:
port: "8080"  # 服务端口

database:
host: "localhost"
user: "your_db_user"
password: "your_db_password"
dbname: "your_database"

key:
news_apikey: "从聚合数据申请的API密钥"
jwt_secret: "生成随机JWT密钥" # 建议长度32位以上

page:
  min_page_size: 1  #最小每页数量
  max_page_size: 20 #最大每页数量
````

### 安装步骤

```bash
# 克隆仓库
git clone https://gitee.com/null_619_3752/app_api.git

# 安装依赖
go mod tidy

