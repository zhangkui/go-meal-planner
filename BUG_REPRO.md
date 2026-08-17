# BUG-004 reproduction

## Task type
diagnosis

## User-visible symptom
库存 500g 时消耗 500g；基线返回 insufficient inventory，精确消耗验证失败。

## Trigger commands
`ash
go test -buildvcs=false -count=1 -run "^TestConsumeExactInventoryLeavesZero$" ./internal/service
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

## Complete baseline evidence
See the private calibration record under the case workspace. The public verification test is present in this branch; hidden answer materials are intentionally excluded.