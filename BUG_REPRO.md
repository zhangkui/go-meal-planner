# BUG-005 reproduction

## Task type
bugfix

## User-visible symptom
查询 2026-08-19 所在周，同时存在 2026-08-17 和 2026-08-24 菜单；基线把下一周周一计入周汇总。

## Trigger commands
`ash
go test -buildvcs=false -count=1 -run "^TestWeeklyNutritionExcludesFollowingMonday$" ./internal/service
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

## Complete baseline evidence
See the private calibration record under the case workspace. The public verification test is present in this branch; hidden answer materials are intentionally excluded.