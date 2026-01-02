# API 文档（草案）

## 统一响应格式
```
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

## 创建短链
- POST /links
- 请求体：
  - url: string
  - expire_at: string(可选，RFC3339，例如 2026-01-02T15:04:05Z)
- 返回：
  - code: string
  - short_url: string

## 重定向
- GET /:code
- 行为：
  - 302 跳转到 original_url

## 获取统计（可选）
- GET /links/:code/stats
- 返回：
  - pv
  - uv
