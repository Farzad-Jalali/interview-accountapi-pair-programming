#!/bin/bash
set -e

testDirs=`go list ./... | grep -v vendor/`

for testDir in $testDirs
do
	go test  -v $testDir
done
