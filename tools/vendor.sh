#!/bin/bash

set -eu

echo "--- Get libs"
go get -u github.com/Sirupsen/logrus
go get -u github.com/husobee/vestigo
go get -u github.com/juju/errors
go get -u github.com/spf13/cobra
go get -u gopkg.in/dancannon/gorethink.v2

echo "--- Vendor libs"
rm -rf vendor
vendor github.com/ernestoalejo/tfg-fn/cmd/fnapi
vendor github.com/ernestoalejo/tfg-fn/cmd/fnctl
