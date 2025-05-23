#!/bin/bash

ZIP=$1
TMPDIR=$(mktemp -d)

# Распаковка
unzip -q "$ZIP" -d "$TMPDIR"

# Редактирование
code "$TMPDIR"

# Ждём подтверждения
echo "Нажми Enter, когда закончишь редактирование и сохранишь все изменения."
read

# Пересборка ZIP
cd "$TMPDIR"
zip -qr "$ZIP" *

# Удаление временной папки
cd -
rm -rf "$TMPDIR"

echo "✅ Обновлённый ZIP сохранён в: $ZIP"

