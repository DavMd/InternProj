#!/bin/bash

# Остановить выполнение скрипта при ошибке
set -e

# Запуск тестов
go test ./... -v | tee results/results.xml

echo "Тестирование завершено"