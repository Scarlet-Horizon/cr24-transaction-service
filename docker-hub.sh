#!/usr/bin/bash

if [ $# -eq 0 ]; then
  echo "Specify tag number!"
  exit 1
fi

if ! [[ $1 =~ ^[+-]?[0-9]+\.?[0-9]*$ ]]; then
  echo "Argument must be a number"
  exit 1
fi

docker compose build
docker tag cr24-transaction-service-api:latest davidslatinek/transaction-api:"$1"
docker push davidslatinek/transaction-api:"$1"
