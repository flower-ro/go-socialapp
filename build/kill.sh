#!/bin/bash

processId=$( netstat -anp | grep socialserver | awk {'split($7, arr, "/"); print arr[1]'}|head -1)

if [ -z "$processId" ]
then
       ps -ef |grep socialserver
	     echo "\$processId is empty"
else
	     echo "\$processId is NOT empty:"$processId
	     kill  -9 $processId
	     processId=""
fi