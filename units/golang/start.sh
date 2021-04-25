#!/bin/bash
make
./golang --config=config.conf


curl -i --request POST 'http://172.35.1.3:9010/testPerf' \
--header 'Content-Type: application/json' \
--data '{
    "lang": "golang",
    "servType": "测试容器",
    "threadCnt": 5,
    "cores": 16,
    "connCnt": 1,
    "reqFreq": 5,
    "keepAlive": 60,
    "packageLen": 0
}'