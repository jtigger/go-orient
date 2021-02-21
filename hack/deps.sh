#!/usr/bin/env bash

set -ex -o pipefail

BASE_DIR=$1

THIS_MOD="$(cat ${BASE_DIR}/go.mod | grep "module" | sed 's|module ||')/" 
# Example: github.com/k14s/kbld

ALL_GO_FILES="${BASE_DIR} -name *.go"

NOT_TESTS_OR_GENERATED="-v \(generated\|_test\.go\)"
# ==>
# ./pkg/kbld/cmd/resolve.go
# ./pkg/kbld/cmd/resolve.go
# ...

EACH_WITH_IMPORTS="grep ${THIS_MOD}"
# ==>
# ./pkg/kbld/cmd/resolve.go:	ctlconf "github.com/k14s/kbld/pkg/kbld/config"
# ./pkg/kbld/cmd/resolve.go:	ctllog "github.com/k14s/kbld/pkg/kbld/logger"
# ...

MAP_FROM_PKG_TO_DEP="s|${BASE_DIR}/\(.*\)/[[:alpha:][:digit:]_]*\.go:[^\"]*\"${THIS_MOD}\([^\"]*\)\".*|\1,\2|"
# ==>
# pkg/kbld/cmd pkg/kbld/config
# pkg/kbld/cmd pkg/kbld/image
# ...

find $ALL_GO_FILES \
  | grep $NOT_TESTS_OR_GENERATED \
  | xargs $EACH_WITH_IMPORTS \
  | sed $MAP_FROM_PKG_TO_DEP \
  | sort --unique

