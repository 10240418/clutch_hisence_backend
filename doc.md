# Database

**Supplier**

- id (uint) - Primary Key
- name (char[64])
- sap (char[16])
- type (enum: 'OEM', 'ODM')

**ProductModel**

- id (uint) - Primary Key
- sap (char[16])
- description (char[128])
- supplier_id (foreignKey) - References Supplier

**ProductionPlan**

- id (uint) - Primary Key
- startAt (date)
- endAt (date)
- belongsTo (char[8])
- product_model_id (foreignKey) - References ProductModel
- planned (int)
- actual (int, nullable)

**ProductLine**

- id (uint) - Primary Key
- name (char[64]) - 可以为空，注册时填写
- pallet_sn_prefix (char[16]) - 可以为空，注册时填写
- device_id (char[64], unique) - 设备唯一标识
- is_registered (bool, default: false) - 是否已注册
- public_key (text) - 注册时生成的公钥

**Pallet**

- id (uint) - Primary Key
- sn (char[32])
- product_model_id (foreignKey) - References ProductModel
- product_line_id (foreignKey) - References ProductLine
- createdAt (dateTime)

**Product**

- id (uint) - Primary Key
- sn (char[32])
- product_model_id (foreignKey) - References ProductModel
- product_line_id (foreignKey) - References ProductLine
- production_plan_id (foreignKey) - References ProductionPlan
- pallet_id (foreignKey) - References Pallet
- createdAt (dateTime)

**User**

- id (int64) - Primary Key
- username (char[32])
- email (string, unique)
- mobile (string, unique)
- password (string, hashed)
- active (bool)
- createdAt (dateTime)
- updatedAt (dateTime)
- deletedAt (dateTime, nullable)

**API**

- id (uint) - Primary Key
- name (char[64])
- app_id (char[32])
- secret (string)

# RESTFul APIs

## Production (产线端接口)

| API                      | Method | Endpoint                       | Auth Required | Description                                  |
| ------------------------ | ------ | ------------------------------ | ------------- | -------------------------------------------- |
| Register Production Line | POST   | `/api/production/register`     | No            | 产线注册，提供 DeviceID 和产线信息获取公钥   |
| Authenticate Device      | POST   | `/api/production/authenticate` | No            | 产线认证，提供 DeviceID 和公钥获取 JWT token |
| Add ProductLine          | POST   | `/api/production/product_line` | ProductLine   | 创建新的生产线                               |
| Delete ProductLine       | DELETE | `/api/production/product_line` | ProductLine   | 删除已有生产线                               |
| Add Pallet               | POST   | `/api/production/pallet`       | ProductLine   | 创建新托盘                                   |
| Add Product              | POST   | `/api/production/product`      | ProductLine   | 创建新产品                                   |

## Management (管理端接口)

| API                   | Method | Endpoint                              | Auth Required | Description               |
| --------------------- | ------ | ------------------------------------- | ------------- | ------------------------- |
| Login                 | POST   | `/api/management/login`               | No            | 管理员登录认证            |
| Add Supplier          | POST   | `/api/management/supplier`            | Admin         | 创建新供应商              |
| Delete Supplier       | DELETE | `/api/management/supplier`            | Admin         | 删除已有供应商            |
| Get Suppliers         | GET    | `/api/management/supplier`            | Admin         | 获取所有供应商列表        |
| Get Supplier          | GET    | `/api/management/supplier/:id`        | Admin         | 获取指定供应商详情        |
| Update Supplier       | PUT    | `/api/management/supplier`            | Admin         | 更新已有供应商            |
| Add ProductModel      | POST   | `/api/management/product_model`       | Admin         | 创建新产品型号            |
| Delete ProductModel   | DELETE | `/api/management/product_model`       | Admin         | 删除已有产品型号          |
| Get ProductModels     | GET    | `/api/management/product_model`       | Admin         | 获取所有产品型号列表      |
| Get ProductModel      | GET    | `/api/management/product_model/:id`   | Admin         | 获取指定产品型号详情      |
| Update ProductModel   | PUT    | `/api/management/product_model`       | Admin         | 更新已有产品型号          |
| Add ProductionPlan    | POST   | `/api/management/production_plan`     | Admin         | 创建新生产计划            |
| Delete ProductionPlan | DELETE | `/api/management/production_plan`     | Admin         | 删除已有生产计划          |
| Get ProductionPlans   | GET    | `/api/management/production_plan`     | Admin         | 获取所有生产计划列表      |
| Get ProductionPlan    | GET    | `/api/management/production_plan/:id` | Admin         | 获取指定生产计划详情      |
| Update ProductionPlan | PUT    | `/api/management/production_plan`     | Admin         | 更新已有生产计划          |
| Add ProductLine       | POST   | `/api/management/product_line`        | Admin         | 录入新产线（仅 DeviceID） |
| Get ProductLines      | GET    | `/api/management/product_line`        | Admin         | 获取所有产线列表          |
| Get ProductLine       | GET    | `/api/management/product_line/:id`    | Admin         | 获取指定产线详情          |
| Delete ProductLine    | DELETE | `/api/management/product_line`        | Admin         | 删除已有产线              |
| Get Pallets           | GET    | `/api/management/pallet`              | Admin         | 获取所有托盘列表          |
| Get Pallet            | GET    | `/api/management/pallet/:id`          | Admin         | 获取指定托盘详情          |
| Get Products          | GET    | `/api/management/product`             | Admin         | 获取所有产品列表          |
| Get Product           | GET    | `/api/management/product/:id`         | Admin         | 获取指定产品详情          |
| Add API               | POST   | `/api/management/api`                 | Admin         | 创建新 API 访问权限       |
| Delete API            | DELETE | `/api/management/api`                 | Admin         | 删除已有 API 访问权限     |
| Get APIs              | GET    | `/api/management/api`                 | Admin         | 获取所有 API 列表         |
| Get API               | GET    | `/api/management/api/:id`             | Admin         | 获取指定 API 详情         |
| Update API            | PUT    | `/api/management/api`                 | Admin         | 更新已有 API 访问权限     |
| Add User              | POST   | `/api/management/user`                | Admin         | 创建新用户                |
| Delete User           | DELETE | `/api/management/user`                | Admin         | 删除已有用户              |
| Get Users             | GET    | `/api/management/user`                | Admin         | 获取所有用户列表          |
| Get User              | GET    | `/api/management/user/:id`            | Admin         | 获取指定用户详情          |
| Update User           | PUT    | `/api/management/user`                | Admin         | 更新已有用户              |

## 设备注册流程

### 1. 管理员录入产线

```json
POST /api/management/product_line
{
  "deviceId": "DEVICE-LINE-001"
}
```

### 2. 产线设备注册

```json
POST /api/production/register
{
  "deviceId": "DEVICE-LINE-001",
  "name": "生产线A",
  "palletSnPrefix": "PLA"
}
```

响应:

```json
{
  "message": "registration successful",
  "publicKey": "generated_public_key_hash"
}
```

### 3. 产线设备认证

```json
POST /api/production/authenticate
{
  "deviceId": "DEVICE-LINE-001",
  "publicKey": "generated_public_key_hash"
}
```

响应:

```json
{
  "message": "authentication successful",
  "token": "jwt_token_for_production_line"
}
```

### 4. 使用 Token 进行后续请求

```
Authorization: Bearer jwt_token_for_production_line
```

## JWT Token 类型

系统支持三种类型的 JWT Token：

1. **Admin Token** (`JwtServiceRoleAdmin`)

   - 用于管理端所有操作
   - 通过管理员登录获取
   - 包含用户 ID 和管理员角色信息

2. **Factory Token** (`JwtServiceRoleFactory`)

   - 预留给工厂级别操作
   - 当前版本未使用

3. **Production Line Token** (`JwtServiceRoleProductionLine`)
   - 用于产线端操作
   - 通过设备认证获取
   - 包含 DeviceID 和产线角色信息

# Frontend Pages

## 1. Authentication

- **Login Page**
  - Login form with username and password

## 2. Supplier Management

- **Supplier List Page**
  - Search/filter functionality
  - Table with columns: ID, Name, SAP, Type
  - Action buttons: View, Edit, Delete
- **Supplier Detail Dialog**
  - Display all supplier details
  - Action buttons: Edit, Delete, Close
- **Add Supplier Dialog**
  - Form with fields: Name, SAP, Type
  - Action buttons: Submit, Cancel
- **Edit Supplier Dialog**
  - Pre-filled form with current values
  - Action buttons: Save, Cancel

## 3. Product Model Management

- **Product Model List Page**
  - Search/filter functionality
  - Table with columns: ID, SAP, Description, Supplier
  - Action buttons: View, Edit, Delete
- **Product Model Detail Dialog**
  - Display all product model details including supplier information
  - Action buttons: Edit, Delete, Close
- **Add Product Model Dialog**
  - Form with fields: SAP, Description, Supplier (dropdown)
  - Action buttons: Submit, Cancel
- **Edit Product Model Dialog**
  - Pre-filled form with current values
  - Action buttons: Save, Cancel

## 4. Production Planning

- **Production Plan List Page**
  - Search/filter functionality
  - Table with columns: ID, StartAt, EndAt, BelongsTo, Product Model, Planned, Actual
  - Action buttons: View, Edit, Delete
- **Production Plan Detail Dialog**
  - Display all production plan details including product model information
  - Action buttons: Edit, Delete, Close
- **Add Production Plan Dialog**
  - Form with fields: StartAt, EndAt, BelongsTo, Product Model (dropdown), Planned, Actual
  - Action buttons: Submit, Cancel
- **Edit Production Plan Dialog**
  - Pre-filled form with current values
  - Action buttons: Save, Cancel

## 5. Production Line Management (产线管理)

- **Product Line List Page**
  - Search/filter functionality
  - Table with columns: ID, DeviceID, Name, PalletSnPrefix, IsRegistered, CreatedAt
  - Action buttons: View, Edit, Delete
  - Registration status indicator
- **Product Line Detail Dialog**
  - Display all product line details including registration status
  - Show public key if registered (masked for security)
  - Action buttons: Edit, Delete, Close
- **Add Product Line Dialog (管理端录入)**
  - Form with fields: DeviceID
  - Action buttons: Submit, Cancel
  - Note: Name and PalletSnPrefix will be filled during device registration
- **Edit Product Line Dialog**
  - Pre-filled form with current values
  - DeviceID cannot be changed if already registered
  - Action buttons: Save, Cancel

## 6. Device Registration Management (设备注册管理)

- **Device Registration Dashboard**
  - Overview of registered/unregistered devices
  - Recent registration activities
  - Device status monitoring
- **Registration Logs**
  - List of all registration attempts
  - Success/failure status
  - Timestamp and device information

## 7. Pallet Management

- **Pallet List Page**
  - Search/filter functionality
  - Table with columns: ID, SN, Product Model, Product Line, CreatedAt
  - Action button: View
- **Pallet Detail Dialog**
  - Display all pallet details including product model and product line information
  - Action button: Close

## 8. Product Management

- **Product List Page**
  - Search/filter functionality
  - Table with columns: ID, SN, Product Model, Product Line, Production Plan, Pallet, CreatedAt
  - Action button: View
- **Product Detail Dialog**
  - Display all product details including related entities
  - Action button: Close

## 9. API Management

- **API List Page**
  - Search/filter functionality
  - Table with columns: ID, Name, AppID
  - Action buttons: View, Edit, Delete
- **API Detail Dialog**
  - Display all API details (excluding Secret)
  - Action buttons: Edit, Delete, Close
- **Add API Dialog**
  - Form with fields: Name, AppID, Secret
  - Action buttons: Submit, Cancel
- **Edit API Dialog**
  - Pre-filled form with current values (excluding Secret)
  - Action buttons: Save, Cancel

## 10. User Management

- **User List Page**
  - Search/filter functionality
  - Table with columns: ID, Username, Email, Mobile, Active, CreatedAt
  - Action buttons: View, Edit, Delete
- **User Detail Dialog**
  - Display all user details (excluding Password)
  - Action buttons: Edit, Delete, Close
- **Add User Dialog**
  - Form with fields: Username, Email, Mobile, Password, Active
  - Action buttons: Submit, Cancel
- **Edit User Dialog**
  - Pre-filled form with current values (excluding Password)
  - Option to reset password
  - Action buttons: Save, Cancel

## 11. Production Operations (产线端操作)

- **Device Registration Interface (产线端注册)**

  - Form with fields: DeviceID, Name, Pallet SN Prefix
  - Action buttons: Register, Cancel
  - Display generated public key upon successful registration
  - Save public key securely for authentication

- **Device Authentication Interface (产线端认证)**

  - Form with fields: DeviceID, Public Key
  - Action buttons: Authenticate, Cancel
  - Display JWT token upon successful authentication
  - Store token for subsequent API calls

- **Production Dashboard**

  - Overview of current production status
  - Quick access to production operations
  - Device registration status

- **Add Product Line Form**

  - Form with fields: Name, Pallet SN Prefix, DeviceID
  - Action buttons: Submit, Cancel
  - Success message and redirect to product line list

- **Add Pallet Form**

  - Form with fields: SN, Product Model (dropdown), Product Line (dropdown)
  - Action buttons: Submit, Cancel
  - Auto-generate SN based on product line prefix
  - Success message and redirect to pallet list

- **Add Product Form**
  - Form with fields: SN, Product Model (dropdown), Product Line (dropdown), Production Plan (dropdown), Pallet (dropdown)
  - Action buttons: Submit, Cancel
  - Cascading dropdowns (Product Line affects available Pallets)
  - Success message and redirect to product list

## Navigation Structure

- **Main Menu**
  - Dashboard
  - Device Management (设备管理)
    - Device Registration (设备注册)
    - Registration Status (注册状态)
  - Production Management
    - Production Plans
    - Production Lines
    - Pallets
    - Products
  - Master Data
    - Suppliers
    - Product Models
  - System Management
    - Users
    - API Access
  - Production Operations (产线端)
    - Device Registration (设备注册)
    - Device Authentication (设备认证)
    - Add Product Line
    - Add Pallet
    - Add Product

## Common UI Components

- **Data Table Component**

  - Pagination support
  - Search/filter functionality
  - Sortable columns
  - Action buttons (View, Edit, Delete)
  - Bulk selection for delete operations

- **Modal Dialog Component**

  - Detail view modal
  - Form modal for add/edit operations
  - Confirmation dialog for delete operations

- **Form Components**
  - Input fields with validation
  - Dropdown selectors with search
  - Date pickers
  - Toggle switches for boolean fields
  - File upload components (if needed)

## Authentication & Authorization

- **Login Required**

  - All pages except login and device registration require authentication
  - JWT token-based authentication
  - Auto-redirect to login if token expired

- **Role-based Access Control**

  - **Admin Role**: All management APIs require admin role
  - **Production Line Role**: Production operations require production line token
  - **Device Registration**: No authentication required for registration and initial authentication
  - Role-based menu visibility

- **Device Registration Security**
  - DeviceID must be pre-registered by admin
  - Public key generated using SHA256 hash of DeviceID
  - JWT token contains device-specific claims
  - Token validation includes role and device identification
