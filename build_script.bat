@echo off
REM Остановить выполнение скрипта при ошибке
setlocal enabledelayedexpansion
REM Создание директории для сборки
if not exist "build" mkdir build
REM Установить зависимости
go mod tidy
REM Сборка проекта
go build -o build ./...
echo Сборка завершена
