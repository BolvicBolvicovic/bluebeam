#!/bin/bash

while ! nc -z mariadb 3306; do
	sleep 0.1
done

go build -o app ./cmd/main.go

./app
