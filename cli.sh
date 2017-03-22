#!/bin/bash

echo $1

curl -H "X-Broker-API-Version: 2.9" localhost:1338/v2/$1




