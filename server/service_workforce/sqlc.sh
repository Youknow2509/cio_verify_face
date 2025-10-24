#!/bin/bash

echo "Generating SQLC code..."
export PATH=$PATH:$(go env GOPATH)/bin
sqlc -f ./config/sqlc.yaml generate