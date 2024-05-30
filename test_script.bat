@echo off
chcp 65001
setlocal enabledelayedexpansion

if not exist "results" mkdir results

go test ./... -v > results/results.xml
type results/results.xml

echo Test Completed
exit /b %errorlevel%
