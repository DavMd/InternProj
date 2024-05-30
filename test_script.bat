@echo off
setlocal enabledelayedexpansion
go test ./graph -v > results/results.xml
type results/results.xml
echo Test Completed
