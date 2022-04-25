#!/bin/bash
workload_path=/root/go/src/github.com/hyperledger/caliper-benchmarks/workload
sign_r=`./test $1 |sed -n '2p'|awk -F " " '{print $2}'`
sign_s=`./test $1 |sed -n '3p'|awk -F " " '{print $2}'`


old_tokenID=`cat $workload_path/queryGrantLog.js|grep contractArguments|awk -F  ":" '{print $2}' |awk -F \" '{print $2}'


old_signR=``
old_signS=``



sed -i s/$old_tokenID/$1/g $workload_path/queryGrantLog.js











