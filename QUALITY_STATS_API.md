# 质量统计 API 文档

## 接口说明

### 获取质量统计数据

**接口地址：** `GET /api/management/quality_stats`

**描述：** 获取指定时间范围内的产品质量统计数据，包括合格率、不良类型分布、供应商不良趋势和各类型不良趋势。

**请求头：**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求参数：**
| 参数名 | 类型 | 必填 | 描述 | 示例 |
|--------|------|------|------|------|
| startDate | string | 是 | 开始日期，格式：YYYY-MM-DD | 2024-01-01 |
| endDate | string | 是 | 结束日期，格式：YYYY-MM-DD | 2024-12-31 |

**请求示例：**

```
GET /api/management/quality_stats?startDate=2024-01-01&endDate=2024-12-31
```

**响应数据结构：**

```json
{
  "data": {
    "qualityRate": {
      "qualifiedCount": 9550,
      "defectCount": 450,
      "totalCount": 10000,
      "qualityRate": 95.5
    },
    "defectTypeDistribution": [
      {
        "type": "绝缘耐压",
        "count": 157,
        "rate": 35.0
      },
      {
        "type": "电阻不良",
        "count": 112,
        "rate": 25.0
      },
      {
        "type": "反电动势",
        "count": 90,
        "rate": 20.0
      },
      {
        "type": "外观",
        "count": 67,
        "rate": 15.0
      },
      {
        "type": "噪音",
        "count": 24,
        "rate": 5.0
      }
    ],
    "supplierDefectTrend": [
      {
        "supplierName": "厂家A",
        "dailyData": [
          {
            "date": "2024-01-01",
            "defectRate": 5.2,
            "totalCount": 500,
            "defectCount": 26
          },
          {
            "date": "2024-01-02",
            "defectRate": 4.8,
            "totalCount": 520,
            "defectCount": 25
          }
        ]
      }
    ],
    "defectTrendByType": {
      "insulationData": [
        {
          "date": "2024-01-01",
          "count": 12
        },
        {
          "date": "2024-01-02",
          "count": 8
        }
      ],
      "resistanceData": [
        {
          "date": "2024-01-01",
          "count": 8
        },
        {
          "date": "2024-01-02",
          "count": 6
        }
      ],
      "emfData": [
        {
          "date": "2024-01-01",
          "count": 4
        },
        {
          "date": "2024-01-02",
          "count": 7
        }
      ],
      "appearanceData": [
        {
          "date": "2024-01-01",
          "count": 2
        },
        {
          "date": "2024-01-02",
          "count": 3
        }
      ],
      "noiseData": [
        {
          "date": "2024-01-01",
          "count": 0
        },
        {
          "date": "2024-01-02",
          "count": 1
        }
      ]
    }
  },
  "message": "success"
}
```

**响应字段说明：**

### qualityRate（合格率统计）

| 字段名         | 类型    | 描述             |
| -------------- | ------- | ---------------- |
| qualifiedCount | int     | 合格产品数量     |
| defectCount    | int     | 不合格产品数量   |
| totalCount     | int     | 总产品数量       |
| qualityRate    | float64 | 合格率（百分比） |

### defectTypeDistribution（不良类型分布）

| 字段名 | 类型    | 描述                               |
| ------ | ------- | ---------------------------------- |
| type   | string  | 不良类型名称                       |
| count  | int     | 该类型不良数量                     |
| rate   | float64 | 该类型在所有不良中的占比（百分比） |

### supplierDefectTrend（供应商不良趋势）

| 字段名       | 类型   | 描述         |
| ------------ | ------ | ------------ |
| supplierName | string | 供应商名称   |
| dailyData    | array  | 每日数据数组 |

#### dailyData 每日数据

| 字段名      | 类型    | 描述                         |
| ----------- | ------- | ---------------------------- |
| date        | string  | 日期（YYYY-MM-DD）           |
| defectRate  | float64 | 该供应商当日不良率（百分比） |
| totalCount  | int     | 该供应商当日总产品数         |
| defectCount | int     | 该供应商当日不良产品数       |

### defectTrendByType（各类型不良趋势）

| 字段名         | 类型  | 描述                 |
| -------------- | ----- | -------------------- |
| insulationData | array | 绝缘耐压不良每日数据 |
| resistanceData | array | 电阻不良每日数据     |
| emfData        | array | 反电动势不良每日数据 |
| appearanceData | array | 外观不良每日数据     |
| noiseData      | array | 噪音不良每日数据     |

#### 每日不良数据结构

| 字段名 | 类型   | 描述               |
| ------ | ------ | ------------------ |
| date   | string | 日期（YYYY-MM-DD） |
| count  | int    | 该类型不良当日数量 |

**错误响应：**

```json
{
  "error": "Invalid start date format, expected YYYY-MM-DD"
}
```

**可能的错误情况：**

- 400: 请求参数格式错误
- 401: 未授权访问
- 500: 服务器内部错误

## 前端集成示例

### 更新合格率饼图

```javascript
// 使用API返回的数据更新合格率饼图
qualityRateOption.value.series[0].data = [
  {
    value: stats.qualityRate.qualityRate,
    name: "合格",
    itemStyle: { color: "#52c41a" },
  },
  {
    value: 100 - stats.qualityRate.qualityRate,
    name: "不合格",
    itemStyle: { color: "#ff4d4f" },
  },
];
```

### 更新不良类型分布饼图

```javascript
defectTypeOption.value.series[0].data = stats.defectTypeDistribution.map(
  (item, index) => ({
    value: item.rate,
    name: item.type,
    itemStyle: { color: getColorByIndex(index) },
  })
);
```

### 更新供应商不良趋势折线图

```javascript
manufacturerTrendOption.value.series = stats.supplierDefectTrend.map(
  (supplier, index) => ({
    name: supplier.supplierName,
    type: "line",
    data: supplier.dailyData.map((day) => day.defectRate),
    itemStyle: { color: getColorByIndex(index) },
  })
);
```

## 数据库查询逻辑

该 API 主要查询以下数据：

1. **合格率计算**：基于 `Product` 表的 `hasDefect` 字段统计
2. **不良类型分布**：基于 `Product` 表的 `defectReason` 字段分组统计
3. **供应商不良趋势**：通过 `Product` → `ProductModel` → `Supplier` 关联查询
4. **各类型不良趋势**：按日期分组统计不同 `defectReason` 的数量

所有查询都会根据 `created_at` 字段过滤指定的时间范围。
