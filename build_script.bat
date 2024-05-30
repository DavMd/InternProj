@echo off
setlocal enabledelayedexpansion
if not exist "build" mkdir build
if not exist "results" mkdir results
go mod tidy
go build -o build ./...
echo Build Completed
