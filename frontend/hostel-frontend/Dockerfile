FROM node:18

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY package*.json ./
RUN npm install --legacy-peer-deps

# Копируем остальные файлы
COPY . .

# Указываем порт
EXPOSE 3001

# Команда запуска
CMD ["npm", "start"]
