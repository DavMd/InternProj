@echo off
setlocal enabledelayedexpansion
go test ./... -v > results/results.xml
type results/results.xml
echo Test Completed
