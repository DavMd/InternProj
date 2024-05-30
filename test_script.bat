@echo off
REM Остановить выполнение скрипта при ошибке
setlocal enabledelayedexpansion
REM Запуск тестов
go test ./... -v | tee results/results.xml
echo Тестирование завершено
