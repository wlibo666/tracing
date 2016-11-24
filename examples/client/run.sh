#!/bin/bash
MAXCNT=10
RUN_CNT=0

while [ $RUN_CNT -lt $MAXCNT ]
do
    curl http://127.0.0.1:1200/img
    RUN_CNT=`expr $RUN_CNT + 1`
done

