#!/bin/bash

for dir in `find . -maxdepth 1 -mindepth 1 -type d`
do
    dir=${dir##*/}
    echo $dir

    cd $dir
    go run main.go

    cd ..
done