#!/bin/bash
previous=$PWD
cd $(dirname $1)
go build -o $previous/$2