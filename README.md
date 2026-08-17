# go-meal-planner

## 项目说明

go-meal-planner 是一个使用内存存储的家庭膳食规划 Web API。它支持菜谱维护、按日期安排三餐、根据菜单和库存生成购物清单、调整食材库存，以及汇总每日和每周营养信息。

项目仅使用 Go 标准库，不依赖数据库、缓存或其他外部服务。服务重启后数据会重置。

## 标准命令

```bash
go build ./...
go test ./...
go vet ./...
go run ./cmd
```

服务默认监听 `:8080`，可通过环境变量 `ADDR` 修改。

## 使用方式

主要接口：

- `POST /recipes`、`GET /recipes`、`GET|PUT|DELETE /recipes/{id}`：菜谱管理
- `PUT /menus/{date}/{meal}`、`GET /menus?from=YYYY-MM-DD&to=YYYY-MM-DD`：菜单规划
- `GET /shopping?from=YYYY-MM-DD&to=YYYY-MM-DD`、`PATCH /shopping/{name}/purchased`：购物清单
- `GET /inventory`、`POST /inventory/{name}/add|consume|adjust`：库存管理
- `GET /nutrition/date/{date}`、`GET /nutrition/week/{date}`：营养汇总

日期使用 `YYYY-MM-DD`，餐次使用 `breakfast`、`lunch` 或 `dinner`。食材单位在同一种食材的菜谱和库存中应保持一致。
