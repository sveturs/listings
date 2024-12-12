FROM node:18

# Устанавливаем рабочую директорию
WORKDIR /app

# Сначала копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install --legacy-peer-deps

# Устанавливаем react-scripts локально в проекте
RUN npm install react-scripts@5.0.1 --save

# Копируем остальные файлы проекта
COPY . .

# Указываем порт
EXPOSE 3001

# Устанавливаем переменные окружения
ENV PORT=3001
ENV WDS_SOCKET_PORT=3001
ENV DANGEROUSLY_DISABLE_HOST_CHECK=true

# Команда запуска
CMD ["npm", "start"]