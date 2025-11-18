# 设备注册功能测试

## API 接口说明

### 1. 产线录入接口（管理端）

- **URL**: `POST /api/management/product_line`
- **描述**: 管理员录入产线 DeviceID
- **认证**: 需要管理员 JWT token
- **请求体**: JSON 格式，只包含 DeviceID

### 2. 产线注册接口

- **URL**: `POST /api/production/register`
- **描述**: 产线端软件通过此接口进行产线注册
- **请求体**: JSON 格式，包含 DeviceID、Name、PalletSnPrefix
- **响应**: 返回生成的公钥

### 3. 产线认证接口

- **URL**: `POST /api/production/authenticate`
- **描述**: 用户通过 DeviceID 和公钥获取 token
- **请求体**: JSON 格式，包含 DeviceID 和 PublicKey
- **响应**: 返回 JWT token

### 4. 需要认证的产线操作接口

- 所有其他产线相关接口都需要在 Header 中携带 JWT token
- **Header**: `Authorization: Bearer <token>`

## 测试流程

1. 管理员先通过管理端录入产线 DeviceID（只需要 DeviceID）
2. 产线端调用注册接口进行注册（提供 DeviceID、Name、PalletSnPrefix）
3. 产线端调用认证接口获取 token（提供 DeviceID 和公钥）
4. 使用 token 调用其他产线操作接口

## 测试数据

见下面的 JSON 文件：

- `01_add_product_line.json` - 管理员录入产线 DeviceID
- `02_register_device.json` - 产线注册（提供完整信息）
- `03_authenticate_device.json` - 产线认证
- `04_authenticated_request.json` - 携带 token 的请求示例
