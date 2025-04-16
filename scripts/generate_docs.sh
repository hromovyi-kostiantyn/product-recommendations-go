#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

MODULE_NAME="product-recommendations-go"
DOCS_DIR="docs/generated"

echo -e "${YELLOW}===== Генерація документації для проєкту =====${NC}"

# Перевірка наявності потрібних інструментів
command -v pkgsite >/dev/null 2>&1 || {
    echo -e "${RED}Помилка: pkgsite не встановлено.${NC}"
    echo -e "Встановіть його командою: go install golang.org/x/pkgsite/cmd/pkgsite@latest"
    exit 1
}

# Створення директорії для документації, якщо вона не існує
mkdir -p ${DOCS_DIR}

echo -e "${YELLOW}Запуск pkgsite для генерації документації...${NC}"

# Створення тимчасової директорії для генерації документації
TEMP_DIR=$(mktemp -d)
echo -e "${YELLOW}Використовуємо тимчасову директорію: ${TEMP_DIR}${NC}"

# Створюємо тимчасовий модуль для генерації документації
cd ${TEMP_DIR}
go mod init example.com/docs
go get ${MODULE_NAME}@latest

# Запускаємо pkgsite в тихому режимі і отримуємо документацію в файлову систему
echo -e "${YELLOW}Генеруємо HTML документацію...${NC}"
mkdir -p site
pkgsite -dir=site -open=false ${MODULE_NAME}

# Копіюємо згенеровану документацію у проєкт
echo -e "${YELLOW}Копіюємо згенеровану документацію...${NC}"
rm -rf ../../${DOCS_DIR}/*
cp -r site/* ../../${DOCS_DIR}/

# Прибираємо тимчасову директорію
cd ../../
rm -rf ${TEMP_DIR}

echo -e "${GREEN}===== Документацію успішно згенеровано в директорії ${DOCS_DIR} =====${NC}"
echo -e "${YELLOW}Для перегляду документації локально виконайте:${NC}"
echo -e "cd ${DOCS_DIR} && python -m http.server 8080"
echo -e "${YELLOW}І відкрийте у браузері http://localhost:8080${NC}"