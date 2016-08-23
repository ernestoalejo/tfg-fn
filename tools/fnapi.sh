#!/bin/bash

set -eu

REGISTRY=$(minikube service --url=true registry | cut -c8-)
VERSION=$(cat /dev/urandom | tr -cd 'a-f0-9' | head -c 8)

echo "--- build app"
go install ./cmd/fnapi

echo "--- stage build"
rm -rf /tmp/build-fnapi
mkdir /tmp/build-fnapi
cp $GOPATH/bin/fnapi /tmp/build-fnapi
cp docker/fnapi/Dockerfile /tmp/build-fnapi

echo "--- build container"
docker build -t $REGISTRY/fnapi:$VERSION /tmp/build-fnapi

echo "--- push container"
docker push $REGISTRY/fnapi:$VERSION

echo "--- patch deployment"
kubectl patch deployment fnapi --patch "$(cat <<EOF
{
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "name": "fnapi",
            "image": "localhost:5000/fnapi:$VERSION"
          }
        ]
      }
    }
  }
}
EOF
)"

echo "--- fnapi successfully deployed to the cluster"
