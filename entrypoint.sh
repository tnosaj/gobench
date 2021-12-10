#!/bin/sh
if [ "unset" = "${DBUSER:-unset}" ]; then
  echo "DBUSER envvar not set"
  exit 1;
fi
if [ "unset" = "${DBPASSWORD:-unset}" ]; then
  echo "DBPASSWORD envvar not set"
  exit 1;
fi
if [ "unset" = "${DBHOSTNAME:-unset}" ]; then
  echo "DBHOSTNAME envvar not set"
  exit 1;
fi
if [ "unset" = "${DBPORT:-unset}" ]; then
  echo "DBPORT envvar not set"
  exit 1;
fi
if [ "unset" = "${DBNAME:-unset}" ]; then
  echo "DBNAME envvar not set"
  exit 1;
fi
if [ "unset" = "${DBENGINE:-unset}" ]; then
  echo "DBENGINE envvar not set"
  exit 1;
fi

# Additional
if [ "unset" = "${ACTION:-unset}" ]; then
  echo "No ACTION: Setting default action of run"
  ACTION="run"
fi
if [ "unset" = "${DURATION:-unset}" ]; then
  echo "No DURATION provided: Setting default of10000000000000 "
  DURATION="10000000000000"
fi
if [ "unset" = "${RATE:-unset}" ]; then
  echo "No RATE provided: disabling rate limiting"
  RATE="0"
fi
if [ "unset" = "${CONNECTIONS:-unset}" ]; then
  echo "No CONNECTIONS provided: Setting default of 10"
  CONNECTIONS="10"
fi
if [ "unset" = "${THREADS:-unset}" ]; then
  echo "No THREADS provided: Setting default of 10"
  THREADS="10"
fi
if [ "unset" = "${SPLIT:-unset}" ]; then
  echo "No SPLIT provided: Setting default of r:w::90:10"
  SPLIT="90:10"
fi
if [ "unset" = "${STRATEGY:-unset}" ]; then
  echo "No STRATEGY provided: Setting default of simple"
  STRATEGY="simple"
fi

# TLS
if [ "unset" = "${CACERT:-unset}" ]; then
  echo "No CACERT provided: Setting default of none"
  CACERT="none"
else
  if [ $(echo "$CACERT" 2>&1 | head -1 | cut -c1-1) != "/" ];then
    echo "$CACERT" > /tmp/ca.crt
    CACERT="/tmp/ca.crt"
  fi
fi
if [ "unset" = "${CLIENTCERT:-unset}" ]; then
  echo "No CLIENTCERT provided: Setting default of none"
  CLIENTCERT="none"
else
  if [ $(echo "$CLIENTCERT" 2>&1 | head -1 | cut -c1-1) != "/" ];then
    echo "$CLIENTCERT" > /tmp/client.crt
    CLIENTCERT="/tmp/client.crt"
  fi
fi
if [ "unset" = "${CLIENTKEY:-unset}" ]; then
  echo "No CLIENTKEY provided: Setting default of none"
  CLIENTKEY="none"
else
  if [ $(echo "$CLIENTKEY" 2>&1 | head -1 | cut -c1-1) != "/" ];then
    echo "$CLIENTKEY" > /tmp/client.key
    chmod 600 /tmp/client.key
    CLIENTKEY="/tmp/client.key"
  fi
fi

/gobench -a "$ACTION" -r "$RATE" -d "$DURATION" -c "$CONNECTIONS" -t "$THREADS" -s "$SPLIT" -cacert "$CACERT" -clientcert "$CLIENTCERT" -clientkey "$CLIENTKEY" -strategy "$STRATEGY"