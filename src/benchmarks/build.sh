#!/bin/bash

for file in $(ls **/*.go); do
    base=$(echo $file | cut -f1 -d.)
    go build -o $base{,.go}
done
