#!/bin/bash
docker-compose build k6
docker-compose up -d app influxdb
for test in load stress spike soak breakpoint; do
  docker-compose run --rm k6 run $test.js
done
