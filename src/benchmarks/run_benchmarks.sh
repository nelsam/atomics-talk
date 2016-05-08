#!/bin/bash

echo Building binaries
./build.sh
echo

run_benchmarks() {
    echo ====================================================
    echo GOMAXPROCS=$GOMAXPROCS
    for file in $(ls *.go); do
        base=$(echo $file | cut -f1 -d.)
        echo Running $base
        time ./$base
        echo
    done
}

export GOMAXPROCS=1
run_benchmarks
export GOMAXPROCS=2
run_benchmarks
export GOMAXPROCS=4
run_benchmarks

echo Cleaning up
./clean.sh
