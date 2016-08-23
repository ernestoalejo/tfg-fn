#!/bin/bash

set -eu

echo "--- start minikube"
minikube start

echo "--- create registry deployment"
kubectl create -f cluster/registry.dep.yml
kubectl create -f cluster/registry.svc.yml
kubectl create -f cluster/registry.ds.yml

echo "--- create RethinkDB deployment"
kubectl create -f cluster/rethinkdb.dep.yml
kubectl create -f cluster/rethinkdb.svc.yml

echo "--- create fnapi deployment"
kubectl create -f cluster/fnapi.dep.yml
kubectl create -f cluster/fnapi.svc.yml

echo "--- cluster created successfully"
