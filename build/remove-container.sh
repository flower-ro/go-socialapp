#!/bin/bash
lastContainerId=$(docker ps -a|grep socialserver-$1 | awk {'print $1'}|head -1)

if [ -z "$lastContainerId" ]
then
       docker ps -all
	     echo "\$lastContainerId is empty"
else
	     echo "\$lastContainerId is NOT empty:"$lastContainerId
	     docker rm -f  $lastContainerId
	     lastContainerId=""
fi