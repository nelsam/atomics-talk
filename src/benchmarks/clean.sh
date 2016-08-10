#!/bin/bash

echo cleaning counting binaries
for file in $(ls **/*.go); do
    base=$(echo $file | cut -f1 -d.)
    rm $base
done
