#!/usr/bin/env bash
# Print the exported-API changes between two git refs of typesense-go.
# Usage: scripts/apidiff.sh <old-ref> <new-ref|WORKTREE> [-incompatible]
# WORKTREE compares against the current checkout, uncommitted changes included.
set -euo pipefail
old=$1 new=$2 flag=${3:-}
repo=$(git rev-parse --show-toplevel)
tmp=$(mktemp -d)
trap 'git -C "$repo" worktree remove --force "$tmp/old" >/dev/null 2>&1 || true; [ "$new" = WORKTREE ] || git -C "$repo" worktree remove --force "$tmp/new" >/dev/null 2>&1 || true; rm -rf "$tmp"' EXIT
GOBIN="$tmp/bin" go install golang.org/x/exp/cmd/apidiff@v0.0.0-20261124000051-d6646c509b6d

export_of() { # <dir> <out>
  local mod
  mod=$(cd "$1" && go list -m)
  (cd "$1" && "$tmp/bin/apidiff" -m -w "$2" "$mod")
}

git -C "$repo" worktree add -q --detach "$tmp/old" "$old"
export_of "$tmp/old" "$tmp/old.exp"
if [ "$new" = WORKTREE ]; then
  export_of "$repo" "$tmp/new.exp"
else
  git -C "$repo" worktree add -q --detach "$tmp/new" "$new"
  export_of "$tmp/new" "$tmp/new.exp"
fi
"$tmp/bin/apidiff" -m $flag "$tmp/old.exp" "$tmp/new.exp"
