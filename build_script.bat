@echo off
setlocal enabledelayedexpansion
if not exist "build" mkdir build
go mod tidy
go build -o build ./...
echo Build Completed
