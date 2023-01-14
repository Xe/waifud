#!/bin/sh

for os in plan9 openbsd darwin; do
	GOOS=$os go build
done
