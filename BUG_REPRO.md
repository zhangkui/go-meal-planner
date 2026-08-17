# BUG-001 reproduction

## Task type
bugfix

## User-visible symptom
提交菜谱更新请求将 steps 设为空数组，随后读取相同菜谱；基线返回旧步骤而不是空数组。完整基线错误见 data/calibration-base-red.txt。

## Trigger commands
`ash
go test -buildvcs=false -count=1 -run "^TestUpdateRecipeCanClearSteps$" ./internal/service
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

## Complete baseline evidence
See the private calibration record under the case workspace. The public verification test is present in this branch; hidden answer materials are intentionally excluded.