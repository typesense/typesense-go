#!/usr/bin/env bash
# Rename the enum constants that typesense-go v4 prefixes with their type name,
# such as api.Upsert to api.IndexActionUpsert. It matches the package by the
# name api and rewrites files in place with gofmt. A second run changes nothing.
# Usage: scripts/v4-enum-rename.sh <path>...
set -euo pipefail
[ "$#" -gt 0 ] || { echo "usage: $0 <path>..." >&2; exit 2; }
while read -r old new; do
  gofmt -r "api.$old -> api.$new" -w "$@"
done <<'EOF'
Counter AnalyticsRuleTypeCounter
Log AnalyticsRuleTypeLog
NohitsQueries AnalyticsRuleTypeNohitsQueries
PopularQueries AnalyticsRuleTypePopularQueries
Contains CurationRuleMatchContains
Exact CurationRuleMatchExact
CoerceOrDrop DirtyValuesCoerceOrDrop
CoerceOrReject DirtyValuesCoerceOrReject
Drop DirtyValuesDrop
Reject DirtyValuesReject
BothSides3 DropTokensModeBothSides3
LeftToRight DropTokensModeLeftToRight
RightToLeft DropTokensModeRightToLeft
Create IndexActionCreate
Emplace IndexActionEmplace
Update IndexActionUpdate
Upsert IndexActionUpsert
EOF
