#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}===== Перевірка якості документації =====${NC}"

# Перевірка наявності потрібних інструментів
command -v golint >/dev/null 2>&1 || {
    echo -e "${RED}Помилка: golint не встановлено.${NC}"
    echo -e "Встановіть його командою: go install golang.org/x/lint/golint@latest"
    exit 1
}

# Перевірка наявності документації для експортованих елементів
echo -e "${YELLOW}Перевірка документації для експортованих елементів...${NC}"
DOC_ISSUES=$(golint -min_confidence=0.3 ./...)

if [ -n "$DOC_ISSUES" ]; then
    echo -e "${RED}Знайдено проблеми з документацією:${NC}"
    echo "$DOC_ISSUES"

    # Підрахунок кількості проблем
    DOC_ISSUES_COUNT=$(echo "$DOC_ISSUES" | grep -c "exported")
    echo -e "${RED}Всього знайдено $DOC_ISSUES_COUNT проблем з документацією експортованих елементів.${NC}"

    # Підрахунок загальної кількості експортованих елементів
    EXPORTED_ELEMENTS=$(grep -r "^func [A-Z]\|^type [A-Z]\|^var [A-Z]\|^const [A-Z]" --include="*.go" . | wc -l)
    echo -e "${YELLOW}Всього експортованих елементів: $EXPORTED_ELEMENTS${NC}"

    # Обчислення відсотка документованих елементів
    DOCUMENTED_PERCENT=$(( 100 - (DOC_ISSUES_COUNT * 100 / EXPORTED_ELEMENTS) ))
    echo -e "${YELLOW}Відсоток документованих елементів: $DOCUMENTED_PERCENT%${NC}"

    if [ $DOCUMENTED_PERCENT -lt 80 ]; then
        echo -e "${RED}Рівень документації нижче 80%. Потрібно покращити документацію!${NC}"
        exit 1
    else
        echo -e "${YELLOW}Рівень документації прийнятний, але можна покращити.${NC}"
    fi
else
    echo -e "${GREEN}Проблем з документацією експортованих елементів не знайдено.${NC}"
fi

# Перевірка наявності коментарів для пакетів
echo -e "${YELLOW}Перевірка документації пакетів...${NC}"
PACKAGES_WITHOUT_DOCS=$(find . -name "*.go" | xargs grep -L "^// Package" | grep -v "_test.go" | grep -v "/vendor/" | sort -u)

if [ -n "$PACKAGES_WITHOUT_DOCS" ]; then
    echo -e "${RED}Знайдено файли без документації пакетів:${NC}"
    echo "$PACKAGES_WITHOUT_DOCS"

    # Підрахунок кількості файлів без документації пакетів
    PACKAGES_WITHOUT_DOCS_COUNT=$(echo "$PACKAGES_WITHOUT_DOCS" | wc -l)
    echo -e "${RED}Всього знайдено $PACKAGES_WITHOUT_DOCS_COUNT файлів без документації пакетів.${NC}"

    # Підрахунок загальної кількості файлів Go
    GO_FILES_COUNT=$(find . -name "*.go" | grep -v "_test.go" | grep -v "/vendor/" | wc -l)
    echo -e "${YELLOW}Всього Go файлів: $GO_FILES_COUNT${NC}"

    # Обчислення відсотка файлів з документацією пакетів
    PACKAGES_DOCUMENTED_PERCENT=$(( 100 - (PACKAGES_WITHOUT_DOCS_COUNT * 100 / GO_FILES_COUNT) ))
    echo -e "${YELLOW}Відсоток файлів з документацією пакетів: $PACKAGES_DOCUMENTED_PERCENT%${NC}"

    if [ $PACKAGES_DOCUMENTED_PERCENT -lt 70 ]; then
        echo -e "${RED}Рівень документації пакетів нижче 70%. Потрібно покращити документацію!${NC}"
        exit 1
    else
        echo -e "${YELLOW}Рівень документації пакетів прийнятний, але можна покращити.${NC}"
    fi
else
    echo -e "${GREEN}Всі пакети мають документацію.${NC}"
fi

# Перевірка наявності прикладів використання
echo -e "${YELLOW}Перевірка наявності прикладів використання...${NC}"
EXAMPLES_COUNT=$(find . -name "*_test.go" | xargs grep -l "func Example" | wc -l)

echo -e "${YELLOW}Знайдено $EXAMPLES_COUNT файлів з прикладами використання.${NC}"

if [ $EXAMPLES_COUNT -lt 3 ]; then
    echo -e "${RED}Недостатньо прикладів використання. Потрібно додати більше прикладів!${NC}"
    echo -e "${RED}Рекомендується створити приклади для основних компонентів системи.${NC}"
else
    echo -e "${GREEN}Кількість прикладів використання прийнятна.${NC}"
fi

echo -e "${GREEN}===== Перевірку якості документації завершено =====${NC}"