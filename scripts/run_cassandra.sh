#!/bin/bash

NAME=archai-cassandra
IMAGE=cassandra:3

PORTS="9042 9160"
PUBLISH=""
for p in $PORTS
do
    PUBLISH="$PUBLISH --publish 127.0.0.1:$p:$p"
done

START="docker start $NAME"
RUN="docker run -d $PUBLISH --name $NAME $IMAGE"
KILL="docker kill $NAME"
RM="docker rm $NAME"
BUILD="docker build -t $NAME ."

echo START=$START
echo RUN=$RUN
echo RM=$RM
echo KILL=$KILL
echo BUILD=$BUILD

TARGET=${1:-default}

case $TARGET in
    start)
        $START
        ;;
    run)
        $RUN
        ;;
    default)
        $START || $RUN
        ;;
    stop)
        $KILL
        $RM
        ;;
    build)
        $BUILD
        ;;
    *)
        echo "$0 (start|run|stop|build)"
        ;;
esac
