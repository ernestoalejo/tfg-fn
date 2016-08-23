#!/bin/bash

set -eu

echo "--- start minikube"
minikube start

echo "--- create registry deployment"
kubectl create -f cluster/registry.dep.yml
kubectl create -f cluster/registry.svc.yml

echo "--- cluster created successfully"
