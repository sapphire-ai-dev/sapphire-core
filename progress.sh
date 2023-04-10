total=$(find . -name '*.go' | sed 's/.*/"&"/' | xargs  wc -l | tail -1 | awk '{ print $1 }')
test=$(find . -name '*test.go' | sed 's/.*/"&"/' | xargs  wc -l | tail -1 | awk '{ print $1 }')
data_json=$(find . -name '*.json' | sed 's/.*/"&"/' | xargs  wc -l | tail -1 | awk '{ print $1 }')
data_sss=$(find . -name '*.sss' | sed 's/.*/"&"/' | xargs  wc -l | tail -1 | awk '{ print $1 }')
data_ssa=$(find . -name '*.ssa' | sed 's/.*/"&"/' | xargs  wc -l | tail -1 | awk '{ print $1 }')
code=$((total-test))
data=$((data_json + data_sss + data_ssa))
total=$((code + test + data))
# shellcheck disable=SC2004
echo "impl: $code, test: $test, data: $data, total: $total"
