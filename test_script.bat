@echo off
setlocal enabledelayedexpansion

echo Checking if results directory exists
if not exist "results" (
    echo Creating results directory
    mkdir results
) else (
    echo Results directory already exists
)

echo Running tests...
go test ./... -v > results/results.xml

echo Test output:
type results/results.xml

echo Test Completed
exit /b %errorlevel%
