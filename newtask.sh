#!/bin/bash

DATA="{ \"task\": \"$1\" }"

echo $DATA

curl -H'Content-Type: application/json' \
  -XPOST http://localhost:8080/addtask -d "$DATA"
