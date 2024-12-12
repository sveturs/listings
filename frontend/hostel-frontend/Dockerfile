FROM node:18

WORKDIR /app

# Копируем файлы проекта
COPY package*.json ./
RUN npm install --legacy-peer-deps

# Установим react-scripts глобально
RUN npm install -g react-scripts

# Копируем остальные файлы
COPY . .

# Указываем порт
EXPOSE 3001

# Команда запуска
CMD ["npm", "start"]