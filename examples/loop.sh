#!/bin/bash

URL="$1"
DB="$2"
BEARER="$3"

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
  curl -H "Authorization:Bearer $BEARER" "$rate" -o "$DB-rate-$action.png"
  curl -H "Authorization:Bearer $BEARER" "$latency" -o "$DB-latency-$action.png"
}

getImages foo bar

exit 0
echo "Looping over: $URL"
fromPrepare=$(date +%s%3N)
curl -s -o /dev/null "$URL/prepare" -X POST -d '{"initialdatasize":50000000}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
toPrepare=$(date +%s%3N)
echo "finished prepare"
echo "starting concurrency 10"
from=10$(date +%s%3N)
curl -s -o /dev/null "$URL/run" -X POST -d '{"concurrency":10,"duration":30000000}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
to10=$(date +%s%3N)
echo "finished concurrency 10"
echo "starting concurrency 20"
from20=$(date +%s%3N)
curl -s -o /dev/null "$URL/run" -X POST -d '{"concurrency":20,"duration":30000000}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
to20=$(date +%s%3N)
echo "finished concurrency 20"
echo "starting concurrency 30"
from30=$(date +%s%3N)
curl -s -o /dev/null "$URL/run" -X POST -d '{"concurrency":30,"duration":30000000}'
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' "$URL/busy")" != "200" ]]; do sleep 5; done
echo "finished concurrency 30"
to30=$(date +%s%3N)

echo "Waiting until prometheus has all the data"
sleep 30

getImages "$fromPrepare" "$toPrepare" prepare
getImages "$from10" "$to10" concurrency10
getImages "$from20" "$to20" concurrency20
getImages "$from30" "$to30" concurrency30


