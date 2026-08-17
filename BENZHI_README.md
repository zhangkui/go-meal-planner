# BENZHI README

## 项目说明
项目：zhangkui/go-meal-planner。这是一个使用 Go 标准库和内存存储的家庭膳食规划 Web API，提供菜谱、菜单、购物清单、库存和营养汇总功能。Go 工具链：golang:1.22；前端工具链：无。

## 标准构建、运行和测试命令
容器内执行：

`ash
cd '/app' && GOTOOLCHAIN=local go build ./...
cd '/app' && GOTOOLCHAIN=local go run ./cmd
cd '/app' && GOTOOLCHAIN=local go test ./...
`

## Docker 构建和进入容器

`ash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh go-meal-planner-bug3-candidate linux/amd64
docker run --rm -it go-meal-planner-bug3-candidate:latest
./build_benzhi_docker.sh go-meal-planner-bug3-candidate-arm64 linux/arm64
docker run --rm -it --platform linux/arm64 go-meal-planner-bug3-candidate-arm64:latest
`

## 题目验证命令

`ash
go test -buildvcs=false -count=1 -run "^TestShoppingListSubtractsInventoryAfterAggregation$" ./internal/service
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

本题类型：bugfix。bugfix candidate 的目标验证命令、全量测试和构建预期退出码为 0；diagnosis candidate 保留基线缺陷，目标复现命令预期为非 0，全量测试反映同一缺陷，构建预期为 0。

## Bug 复现

现象和触发步骤见 BUG_REPRO.md。容器中不包含隐藏测试、修复补丁或提示词。