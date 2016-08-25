#!/bin/bash

run_benchmarks() {
    echo ====================================================
    echo GOMAXPROCS=$GOMAXPROCS
    echo GOROUTINES=$GOROUTINES
    for file in $(ls **/*.go); do
        base=$(echo $file | cut -f1 -d.)
        echo running $base
        time ./$base
        echo
    done
}

main() {
    echo building binaries
    ./build.sh
    echo

    for i in 1 2 4 8 16; do
        export GOMAXPROCS=$i
        for j in 2 4 8 16; do
            export GOROUTINES=$j
            run_benchmarks
        done
    done

    echo cleaning up
    ./clean.sh
}

main
