#!/bin/bash

cd "$(dirname "$0")"

sqlc -f ./config/sqlc.yaml generate
