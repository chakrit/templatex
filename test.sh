#!/bin/sh
go run extend.go | diff - output.txt
