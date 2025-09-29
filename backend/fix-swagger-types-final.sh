#!/bin/bash

echo "Исправляем Swagger аннотации во всех файлах..."

# Найти все Go файлы и заменить utils.SuccessResponse и utils.ErrorResponse
find . -name "*.go" -type f -exec sed -i \
  -e 's/@Success \([0-9]*\) {object} utils\.SuccessResponse/@Success \1 {object} backend_pkg_utils.SuccessResponseSwag/g' \
  -e 's/@Failure \([0-9]*\) {object} utils\.ErrorResponse/@Failure \1 {object} backend_pkg_utils.ErrorResponseSwag/g' \
  {} \;

echo "Готово!"