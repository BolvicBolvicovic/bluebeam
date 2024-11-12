# !/bin/bash

TEMPFILE="/tmp/bluebeam_test_unit.txt"

touch ${TEMPFILE}

go test -c ./...
echo

for TEST in *.test;
do
	./${TEST} -test.v >> ${TEMPFILE} 2>&1
done

grep -E  -- "--- PASS|--- FAIL" ${TEMPFILE}

total_pass=$(grep -E "^--- PASS" "${TEMPFILE}" | wc -l)
total_fail=$(grep -E "^--- FAIL" "${TEMPFILE}" | wc -l)
echo
echo "Total: $total_pass PASS, $total_fail FAIL"

rm ./*.test ${TEMPFILE}
