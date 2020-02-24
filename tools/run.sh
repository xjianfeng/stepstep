#!/bin/bash
cmd=$1
case $cmd in 
    start)
	nohup ./server 2>&1 >> game.log &
	;;
     stop)
        killall server
	;;
     restart)
        killall server
	nohup ./server 2>&1 >> game.log &
	;;
    *)
    ;;
esac
