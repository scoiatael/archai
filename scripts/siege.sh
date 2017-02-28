#!/usr/bin/env bash

siege -t60S -b -H 'Content-Type: application/json' "http://localhost:8080/stream/siege-test POST < $(dirname $0)/siege_data.json"
