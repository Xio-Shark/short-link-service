# 设计说明

## 核心流程
1. 创建短链：校验 URL -> 生成 code -> 入库 -> 缓存
2. 重定向：读取缓存 -> 命中则跳转 -> 未命中查库并回写缓存

## 数据模型（初版）
- links
  - id（自增主键）
  - code（唯一短码）
  - original_url（原始链接）
  - expire_at（过期时间，时间戳）
  - status（1=启用，0=禁用）
  - created_at（创建时间）
  - updated_at（更新时间）

- visits
  - id（自增主键）
  - link_id（关联 links.id）
  - ip（访问IP）
  - user_agent（UA）
  - created_at（访问时间）

## 建表 SQL（最小可用）
详见 `docs/schema.sql`。

## 关键决策
- 短链策略：Base62 + 自增ID，保证短码稳定且递增可控
- 冲突处理：code 唯一索引，写入冲突则重试生成（最多 N 次）
- 缓存策略：缓存 code -> original_url，过期与 links 保持一致
- 统计策略：同步写库或异步批量

## 短链生成策略（含冲突处理）
1. 写入原始链接，获取自增 ID（可先插入空 code）
2. 使用 Base62 编码生成 code
3. 更新 links.code，若唯一索引冲突则重试（重新生成或追加随机盐）
4. 成功后写入缓存：code -> original_url（TTL 与 expire_at 对齐）

建议重试次数 3-5 次，仍失败返回错误码并记录告警日志。
