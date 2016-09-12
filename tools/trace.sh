#!/bin/bash

set -eu

for i in `seq 1 1000`;
do
  fnctl functions list
done 
