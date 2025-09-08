#!/bin/bash

echo "Проверка использования API endpoints в frontend коде..."
echo "=========================================="
echo ""

echo "Файлы, которые всё еще используют относительные пути (fetch с '/api'):"
echo "---------------------------------------------------------------------"
grep -r "fetch(['\"]\/api" frontend/svetu/src --include="*.ts" --include="*.tsx" --include="*.js" --include="*.jsx" 2>/dev/null | grep -v "node_modules" | grep -v ".next" | cut -d: -f1 | sort | uniq

echo ""
echo "Файлы, которые всё еще используют axios с относительными путями:"
echo "---------------------------------------------------------------"
grep -r "axios.*['\"]\/api" frontend/svetu/src --include="*.ts" --include="*.tsx" --include="*.js" --include="*.jsx" 2>/dev/null | grep -v "node_modules" | grep -v ".next" | cut -d: -f1 | sort | uniq

echo ""
echo "=========================================="
echo "Эти файлы нужно обновить для использования configManager.get('api.url')"