#!/bin/bash

echo "Running $1 test cases"
for testName in cnf_tests/$1/*.cnf; do
    echo -- "$testName" ---
    ./run.sh "$testName"
done