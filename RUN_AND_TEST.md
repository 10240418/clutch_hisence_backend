# 运行和测试指南

## 前置要求

1. **Go 环境**: 需要 Go 1.23.0 或更高版本
2. **MySQL 数据库**: 需要运行 MySQL 数据库服务器
3. **环境变量配置**: 需要创建 `.env` 文件

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境变量

复制 `.env.example` 文件为 `.env` 并修改配置：

```bash
# Windows PowerShell
Copy-Item .env.example .env

# Linux/Mac
cp .env.example .env
```

然后编辑 `.env` 文件，修改数据库连接信息：

```env
DB_HOST=localhost
DB_USER=root
DB_PASS=your_password
DB_PORT=3306
DB_NAME=hisense_vmi
HTTP_HOST=0.0.0.0
HTTP_PORT=8080
LOG_FILE=hisense-vmi-server.log
SECRET=your-secret-key-here
PRODUCTION=false
```

### 3. 创建数据库

在 MySQL 中创建数据库：

```sql
CREATE DATABASE hisense_vmi CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 运行服务器

```bash
go run main.go preload.go global.go
```

或者编译后运行：

```bash
go build -o main.exe .
./main.exe
```

服务器启动后，默认运行在 `http://localhost:8080`

## 测试 API

### 方法 1: 使用 curl 命令

#### 1. 管理员登录

```bash
curl -X POST http://localhost:8080/api/management/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"admin@admin.admin\",\"password\":\"admin\"}"
```

响应示例：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 2. 获取供应商列表（需要管理员 token）

```bash
curl -X GET http://localhost:8080/api/management/supplier \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 方法 2: 使用 PowerShell 脚本测试

在 Windows 上可以使用 PowerShell 脚本进行测试。查看 `test_data` 目录中的测试脚本。

### 方法 3: 使用 Postman 或类似工具

1. 导入 API 集合（如果有）
2. 设置基础 URL 为 `http://localhost:8080`
3. 先调用登录接口获取 token
4. 在后续请求的 Header 中添加：`Authorization: Bearer <token>`

## 测试数据报表接口

### 不合格报表查询

```bash
curl -X GET "http://localhost:8080/api/management/report/defect?startDate=2024-01-01&endDate=2024-12-31&pageNum=1&pageSize=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 检测报表查询

```bash
curl -X GET "http://localhost:8080/api/management/report/inspection?startDate=2024-01-01&endDate=2024-12-31&pageNum=1&pageSize=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 检测费用报表查询

```bash
curl -X GET "http://localhost:8080/api/management/report/cost?startDate=2024-01-01&endDate=2024-12-31&pageNum=1&pageSize=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 设备注册流程测试

参考 `test_data/device_registration_tests.md` 文档和 `test_data/test_device_registration.sh` 脚本。

在 Windows 上，可以使用 Git Bash 运行测试脚本，或者使用 PowerShell 版本。

## 常见问题

### 1. 数据库连接失败

- 检查 MySQL 服务是否运行
- 验证 `.env` 文件中的数据库配置是否正确
- 确认数据库用户有足够的权限

### 2. 端口被占用

- 修改 `.env` 文件中的 `HTTP_PORT` 为其他端口（如 8081）
- 或者关闭占用 8080 端口的其他程序

### 3. 日志文件权限错误

- 确保应用有权限在项目目录创建和写入日志文件
- 或者修改 `LOG_FILE` 路径为有权限的目录

### 4. 默认管理员账户

系统会自动创建默认管理员账户：
- 用户名/邮箱: `admin@admin.admin`
- 密码: `admin`
- 手机号: `admin`

**注意**: 生产环境请立即修改默认密码！

## 使用 Docker 运行

### 构建镜像

```bash
docker build -t hisense-vmi-server .
```

### 运行容器

```bash
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=your_db_host \
  -e DB_USER=your_db_user \
  -e DB_PASS=your_db_password \
  -e DB_PORT=3306 \
  -e DB_NAME=hisense_vmi \
  -e HTTP_HOST=0.0.0.0 \
  -e HTTP_PORT=8080 \
  -e LOG_FILE=/app/hisense-vmi-server.log \
  -e SECRET=your-secret-key \
  -e PRODUCTION=true \
  --name hisense-vmi-server \
  hisense-vmi-server
```

## 开发模式

在开发模式下，应用会：
- 自动加载 `.env` 文件
- 自动创建数据库表（通过 migration）
- 自动创建默认管理员账户（如果不存在）

## 生产模式

设置环境变量 `PRODUCTION=true` 时：
- 不会加载 `.env` 文件（使用环境变量）
- 启动前等待 10 秒（等待数据库等依赖服务就绪）



