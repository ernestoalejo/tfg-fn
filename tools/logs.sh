#!/bin/bash

set -eu

while true; do
  kubectl logs -f $(kubectl get pods -lapp=fnapi -o jsonpath="{.items[0].metadata.name}")
  echo "------------------------------------- RESTART"
  sleep 2
done
