#!/bin/bash
# Скрипт для проверки конкретного файла TypeScript

# Получаем путь к файлу из аргумента
FILE=$1

if [ -z "$FILE" ]; then
  echo "Необходимо указать путь к файлу для проверки"
  echo "Пример: ./check-file.sh src/pages/marketplace/ListingDetailsPage.tsx"
  exit 1
fi

# Создаем временную копию файла для проверки
TEMP_DIR="temp-check"
mkdir -p $TEMP_DIR
TEMP_FILE="$TEMP_DIR/$(basename $FILE)"
cp $FILE $TEMP_FILE

# Используем временный tsconfig для проверки
cat > $TEMP_DIR/tsconfig.json << EOF
{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "strict": false,
    "forceConsistentCasingInFileNames": true,
    "module": "esnext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react",
    "noImplicitAny": false
  },
  "include": ["$(basename $FILE)"],
  "exclude": ["../node_modules"]
}
EOF

echo "Проверка файла $FILE..."

# Используем dummy import вместо реальных
sed -i -e 's/import /\/\/ import /g' \
       -e 's/export /\/\/ export /g' $TEMP_FILE

# Проверяем только синтаксис TypeScript, без зависимостей
cd $TEMP_DIR && NODE_OPTIONS="--max-old-space-size=4096" npx tsc --noEmit --allowJs --checkJs false $(basename $FILE)
CHECK_RESULT=$?

# Возвращаемся в исходную директорию и удаляем временные файлы
cd ..
rm -rf $TEMP_DIR

if [ $CHECK_RESULT -eq 0 ]; then
  echo "✅ Файл $FILE не содержит синтаксических ошибок TypeScript"
else
  echo "❌ Файл $FILE содержит ошибки TypeScript"
fi

exit $CHECK_RESULT