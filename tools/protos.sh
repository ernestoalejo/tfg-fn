#!/bin/bash

set -eu

protoc protos/fn.proto --go_out=plugins=grpc:.
