# 短链接服务（Gin + MySQL + Redis）

## 目标
构建一个可对外演示的短链接服务，具备核心功能、基本工程化与性能指标。

## 技术栈
- Go / Gin
- MySQL / GORM
- Redis
- Docker（可选）

## 依赖安装
```
go env -w GOPROXY=https://proxy.golang.org,direct
go mod tidy
```

## 目录结构
```
short-link-service/
  cmd/                # 程序入口
  internal/           # 业务代码
  configs/            # 配置文件
  docs/               # 文档与设计说明
```

## 功能范围
- 生成短链接（Base62 或 Hash）
- 短链接重定向
- 过期时间与禁用
- 访问统计（PV/UV）
- 后台管理（可选）

## 质量目标
- 2000 QPS 下重定向 P95 < 50ms
- Redis 命中率 >= 90%
- 关键接口具备基本日志与错误码

## 文档产出
- `docs/design.md` 设计说明
- `docs/api.md` 接口文档
- `docs/tasks.md` 任务清单
- `docs/deploy.md` 部署建议

## 快速启动（Docker Compose）
```
docker compose up -d
```

## 本地启动
```
APP_CONFIG=configs/config.example.yaml go run ./cmd
```

## 配置说明
- `configs/config.example.yaml` 示例配置
- `configs/config.dev.yaml` Docker Compose 配置
- `configs/config.prod.yaml` 生产环境模板
