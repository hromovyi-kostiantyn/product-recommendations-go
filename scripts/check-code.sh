#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}===== Перевірка коду проєкту =====${NC}"

# Перевірка форматування
echo -e "${YELLOW}Перевірка форматування коду...${NC}"
if gofmt -l . | grep -q .; then
    echo -e "${RED}Знайдено неправильно відформатовані файли:${NC}"
    gofmt -l .
    echo -e "${YELLOW}Виконайте команду 'gofmt -w .' для виправлення.${NC}"
    exit 1
else
    echo -e "${GREEN}Перевірка форматування успішна!${NC}"
fi

# Перевірка імпортів
echo -e "${YELLOW}Перевірка імпортів...${NC}"
if goimports -l . | grep -q .; then
    echo -e "${RED}Знайдено проблеми з імпортами:${NC}"
    goimports -l .
    echo -e "${YELLOW}Виконайте команду 'goimports -w .' для виправлення.${NC}"
    exit 1
else
    echo -e "${GREEN}Перевірка імпортів успішна!${NC}"
fi

# Запуск go vet
echo -e "${YELLOW}Запуск go vet...${NC}"
if ! go vet ./...; then
    echo -e "${RED}go vet виявив проблеми.${NC}"
    exit 1
else
    echo -e "${GREEN}go vet не виявив проблем!${NC}"
fi

# Запуск golangci-lint
echo -e "${YELLOW}Запуск golangci-lint...${NC}"
if ! golangci-lint run ./...; then
    echo -e "${RED}golangci-lint виявив проблеми.${NC}"
    exit 1
else
    echo -e "${GREEN}golangci-lint не виявив проблем!${NC}"
fi

# Запуск тестів
echo -e "${YELLOW}Запуск тестів...${NC}"
if ! go test ./... -cover; then
    echo -e "${RED}Тести не пройшли.${NC}"
    exit 1
else
    echo -e "${GREEN}Тести пройшли успішно!${NC}"
fi

echo -e "${GREEN}===== Всі перевірки пройдено успішно! =====${NC}"