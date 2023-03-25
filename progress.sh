total=$(find . -name '*.go' | sed 's/.*/"&"/' | xargs  wc -l | tail -1 | awk '{ print $1 }')
test=$(find . -name '*test.go' | sed 's/.*/"&"/' | xargs  wc -l | tail -1 | awk '{ print $1 }')
# shellcheck disable=SC2004
echo "impl: $(($total-$test)), test: $test, total: $total"
