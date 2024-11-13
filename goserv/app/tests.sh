#!/bin/bash

TEMPFILE="/tmp/bluebeam_test_unit.txt"
COVERDIR="coverdir"
COVERFILE="coverfile"

touch "${TEMPFILE}"
mkdir -p "${COVERDIR}"

export "GOCOVERDIR=${COVERDIR}"
go test -c -cover ./...
echo

for TEST in *.test; do
  COVERFILE_PATH="${COVERDIR}/${COVERFILE}_${TEST}"

  echo "${TEST} is executed:"
  echo
  ./${TEST} -test.v -test.coverprofile="${COVERFILE_PATH}" | grep -E -- "--- PASS|--- FAIL" | tee -a ${TEMPFILE} 2<&1
  echo
  
done


total_pass=$(grep -E "^--- PASS" "${TEMPFILE}" | wc -l)
total_fail=$(grep -E "^--- FAIL" "${TEMPFILE}" | wc -l)
echo
echo "Total: $total_pass PASS, $total_fail FAIL"
echo

# This command comes from https://github.com/wadey/gocovmerge
gocovmerge ${COVERDIR}/* > ${COVERDIR}/${COVERFILE}

go tool cover -func="${COVERDIR}/${COVERFILE}"

if [ "$1" == "-html" ]; then
  go tool cover -html="${COVERDIR}/${COVERFILE}"
fi

rm -rf ./*.test ${TEMPFILE} ${COVERDIR}
