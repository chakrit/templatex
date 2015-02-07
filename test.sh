#!/bin/sh
go run *.go | diff - output.txt
