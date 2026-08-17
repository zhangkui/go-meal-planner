# BUG-002 reproduction

## Task type
diagnosis

## User-visible symptom
同一日期和餐次先创建未确认菜单，再提交另一菜谱；基线第二次提交成功并覆盖原菜单，验证测试返回 expected conflict but got nil。

## Trigger commands
`ash
go test -buildvcs=false -count=1 -run "^TestUnconfirmedMenuCannotBeOverwritten$" ./internal/service
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

## Complete baseline evidence
See the private calibration record under the case workspace. The public verification test is present in this branch; hidden answer materials are intentionally excluded.