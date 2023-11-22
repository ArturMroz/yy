#!/bin/sh

# set -e

rm -f yy
go build -cover -o yy .

rm -rf covdatafiles
mkdir covdatafiles

errors=0

FILES=$(find examples -name "*.yeet" -print)
for F in $FILES
do
  out=`GOCOVERDIR=covdatafiles ./yy $F`
  if [ $? -ne 0 ]; then
    echo FAIL: $F
    echo "$out"
    errors=$((errors+1))
  fi
done

N=$(echo $FILES | wc -w)

if [ $errors -gt 0 ]; then
  echo "run $N files, $errors errors"
  exit
else
  echo "run $N files, no errors"
fi

# Post-process the resulting profiles.
go tool covdata percent -i=covdatafiles

go tool covdata textfmt -i=covdatafiles -o=cov.txt
go tool cover -html=cov.txt
