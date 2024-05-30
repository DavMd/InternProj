@echo off
setlocal enabledelayedexpansion
go mod tidy
go build -o build ./...
echo Сборка завершена
