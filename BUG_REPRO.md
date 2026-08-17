# BUG-003 reproduction

## Task type
bugfix

## User-visible symptom
两顿菜单各需要 2 个鸡蛋且库存 3 个；基线逐餐扣库存，购物清单漏掉鸡蛋。完整基线错误见 data/calibration-base-red.txt。

## Trigger commands
`ash
go test -buildvcs=false -count=1 -run "^TestShoppingListSubtractsInventoryAfterAggregation$" ./internal/service
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

## Complete baseline evidence
See the private calibration record under the case workspace. The public verification test is present in this branch; hidden answer materials are intentionally excluded.