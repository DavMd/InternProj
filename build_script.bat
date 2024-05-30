@echo off
REM Остановить выполнение скрипта при ошибке
setlocal enabledelayedexpansion
REM Установить зависимости
go mod tidy
REM Сборка проекта
go build -o InternProj ./...
echo Сборка завершена
