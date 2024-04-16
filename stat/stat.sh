#!/bin/sh

prev=
for d in $(ls); do
    if ! [ 20240329${d#20240329} = $d ]; then
        continue
    fi

    for i in 1 2 3; do
        if ! [ -f $d/less$i.csv ] && [ -f $prev/less$i.csv ]; then
            cp $prev/less$i.csv $d/less$i.csv
        fi
    done

    if ! [ -f $d/less4.csv ] && [ -f $d/less4.html ]; then
        ./bin/parse $d/less4.html > $d/less4.csv
    fi

    echo ==$d
    (cd $d && ../bin/stat 4 2>&1 | grep -v .csv | head -4) 
    echo ...
    echo

    prev=$d
done