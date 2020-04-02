#!/usr/bin/env vbash

source "${BASH_SOURCE%/*}/common.bash"

shopt -s extglob

go mod download

tools=$(go list -m -f={{.Dir}} golang.org/x/tools)

echo "Tools is $tools"

cd $(git rev-parse --show-toplevel)
regex='s+golang.org/x/tools/internal+github.com/govim/govim/cmd/govim/internal/golang_org_x_tools+g'

rsync -a --delete --chmod=Du+w,Fu+w $tools/internal/ ./cmd/govim/internal/golang_org_x_tools
find ./cmd/govim/internal/golang_org_x_tools/ -name "*_test.go" -exec rm {} +
find ./cmd/govim/internal/golang_org_x_tools -name "*.go" -exec sed -i'' -e $regex {} +
cp $tools/LICENSE ./cmd/govim/internal/golang_org_x_tools

go mod tidy
