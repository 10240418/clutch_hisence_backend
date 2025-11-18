#!/bin/bash

# 设备注册功能测试脚本
# 请确保服务器在 localhost:8080 运行

BASE_URL="http://localhost:8080"
ADMIN_TOKEN=""

echo "=== 设备注册功能测试 ==="

# 1. 管理员登录获取token
echo "1. 管理员登录..."
ADMIN_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/management/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin@admin.admin",
    "password": "admin"
  }')

echo "管理员登录响应: $ADMIN_LOGIN_RESPONSE"

# 提取token (需要jq工具)
if command -v jq &> /dev/null; then
    ADMIN_TOKEN=$(echo $ADMIN_LOGIN_RESPONSE | jq -r '.token')
    echo "获取到管理员token: ${ADMIN_TOKEN:0:50}..."
else
    echo "请安装jq工具或手动提取token"
    exit 1
fi

# 2. 管理员录入产线信息（只录入DeviceID）
echo -e "\n2. 管理员录入产线DeviceID..."
ADD_PRODUCTLINE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/management/product_line" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d @test_data/01_add_product_line.json)

echo "录入产线DeviceID响应: $ADD_PRODUCTLINE_RESPONSE"

# 3. 产线注册（提供DeviceID、Name、PalletSnPrefix）
echo -e "\n3. 产线设备注册（提供完整信息）..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/production/register" \
  -H "Content-Type: application/json" \
  -d @test_data/02_register_device.json)

echo "设备注册响应: $REGISTER_RESPONSE"

# 提取公钥
if echo $REGISTER_RESPONSE | jq -e '.publicKey' > /dev/null; then
    PUBLIC_KEY=$(echo $REGISTER_RESPONSE | jq -r '.publicKey')
    echo "获取到公钥: $PUBLIC_KEY"
    
    # 更新认证请求文件
    cat > test_data/03_authenticate_device_updated.json <<EOF
{
  "deviceId": "DEVICE-LINE-001",
  "publicKey": "$PUBLIC_KEY"
}
EOF
    
    # 4. 产线认证
    echo -e "\n4. 产线设备认证..."
    AUTH_RESPONSE=$(curl -s -X POST "$BASE_URL/api/production/authenticate" \
      -H "Content-Type: application/json" \
      -d @test_data/03_authenticate_device_updated.json)
    
    echo "设备认证响应: $AUTH_RESPONSE"
    
    # 提取产线token
    if echo $AUTH_RESPONSE | jq -e '.token' > /dev/null; then
        DEVICE_TOKEN=$(echo $AUTH_RESPONSE | jq -r '.token')
        echo "获取到设备token: ${DEVICE_TOKEN:0:50}..."
        
        # 5. 使用token进行认证请求
        echo -e "\n5. 使用token创建托盘..."
        PALLET_RESPONSE=$(curl -s -X POST "$BASE_URL/api/production/pallet" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $DEVICE_TOKEN" \
          -d @test_data/04_authenticated_request.json)
        
        echo "创建托盘响应: $PALLET_RESPONSE"
        
        echo -e "\n=== 测试完成 ==="
    else
        echo "设备认证失败"
    fi
else
    echo "设备注册失败"
fi
