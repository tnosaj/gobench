#!/bin/bash

URL="$1"
echo "Looping over: $URL"
curl -s -o /dev/null "$URL/prepare" -X POST -d '{"initialdatasize":50000000}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
echo "finished prepare"
echo "starting concurrency 10"
curl -s -o /dev/null "$URL/run" -X POST -d '{"concurrency":10,"duration":30000000}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
echo "finished concurrency 10"
echo "starting concurrency 20"
curl -s -o /dev/null "$URL/run" -X POST -d '{"concurrency":20,"duration":30000000}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
echo "finished concurrency 20"
echo "starting concurrency 30"
curl -s -o /dev/null "$URL/run" -X POST -d '{"concurrency":30,"duration":30000000}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
echo "finished concurrency 30"
