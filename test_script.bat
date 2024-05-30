@echo off
setlocal enabledelayedexpansion
go test ./... -v | tee results/results.xml
echo Test Completed
