# 实战部署建议

## 本地开发（快速验证）
1. 启动 MySQL 与 Redis（推荐 Docker）
```
docker run --name shortlink-mysql -e MYSQL_ROOT_PASSWORD=pass -e MYSQL_DATABASE=shortlink -p 3306:3306 -d mysql:8

docker run --name shortlink-redis -p 6379:6379 -d redis:7
```
2. 初始化表结构
```
mysql -h 127.0.0.1 -P 3306 -u root -ppass shortlink < docs/schema.sql
```
3. 启动服务
```
APP_CONFIG=configs/config.example.yaml go run ./cmd
```

## 预发布环境（推荐方式）
- 使用 Docker Compose 统一编排（服务 + MySQL + Redis）
- 配置分环境：`configs/config.dev.yaml` / `configs/config.prod.yaml`
- 开启基础监控：日志收集 + 简单健康检查 `/healthz`

### Docker Compose 示例
```
docker compose up -d
```

### 数据初始化
```
docker exec -i shortlink-mysql mysql -u root -ppass shortlink < docs/schema.sql
```

## 线上环境（最小成本）
- 规格建议：2 vCPU / 4 GB RAM 起步（访问量增大后扩容）
- 系统：Ubuntu 22.04 LTS
- 数据库：优先云数据库（RDS）以降低维护成本
- 缓存：Redis 建议使用托管服务或哨兵模式

## 可复用运维清单
- 日志：访问日志 + 错误日志 + 请求ID
- 监控：服务存活 + QPS + 延迟分位
- 备份：MySQL 定期全量 + binlog
- 安全：最少权限账号、端口限制、HTTPS 代理（Nginx）
