#!/bin/bash

set -eu

echo "GET http://192.168.99.100:31073/trigger/simple?name=Ernesto" | vegeta attack -duration=30s -rate 15 | vegeta report -reporter plot > result.html
