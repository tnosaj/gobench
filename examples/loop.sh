#!/bin/bash

URL="$1"
DB="$2"
BEARER="$3"

run=$(date +%s)
imageDir="$HOME/images/$run"

mkdir -p "$imageDir"

latencyurl="http://localhost:3000/render/d-solo/ddrgi9ifi18n4c/gobench?orgId=1&from=FROMDATE&to=TODATE&panelId=1&width=1000&height=500&scale=1&tz=Europe%2FVienna"
rateurl="http://localhost:3000/render/d-solo/ddrgi9ifi18n4c/gobench?orgId=1&from=FROMDATE&to=TODATE&panelId=2&width=1000&height=500&scale=1&tz=Europe%2FVienna"

getImages () {
  from="$1"
  to="$2"
  action="$3"
  echo "Getting images from $from to $to for $action"
  tmplatency="${latencyurl/FROMDATE/$from}"
  latency="${tmplatency/TODATE/$to}"
  tmprate="${rateurl/FROMDATE/$from}"
  rate="${tmprate/TODATE/$to}"
  curl -s -H "Authorization:Bearer $BEARER" "$rate" -o "$imageDir/$DB-rate-$action.png"
  curl -s -H "Authorization:Bearer $BEARER" "$latency" -o "$imageDir/$DB-latency-$action.png"
}

echo "Looping over: $URL"
fromPrepare=$(date +%s%3N)
curl -s -o /dev/null "$URL/prepare" -X POST -d '{"initialdatasize":2000000,"concurrency":15}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
echo "finished prepare"
echo "Waiting until prometheus has all the data"
sleep 30
toPrepare=$(date +%s%3N)
getImages "$fromPrepare" "$toPrepare" prepare
echo "starting concurrency 10"
from10=$(date +%s%3N)
curl -s -o /dev/null "$URL/run" -X POST -d '{"concurrency":5,"duration":2000000}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
echo "finished concurrency 10"
echo "Waiting until prometheus has all the data"
sleep 30
to10=$(date +%s%3N)
getImages "$from10" "$to10" concurrency5
echo "starting concurrency 20"
from20=$(date +%s%3N)
curl -s -o /dev/null "$URL/run" -X POST -d '{"concurrency":10,"duration":2000000}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
echo "finished concurrency 20"
echo "Waiting until prometheus has all the data"
sleep 30
to20=$(date +%s%3N)
getImages "$from20" "$to20" concurrency10
