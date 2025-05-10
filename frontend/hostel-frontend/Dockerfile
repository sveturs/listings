FROM node:20

WORKDIR /app

# Копируем файлы package.json и package-lock.json
COPY package*.json ./

# Устанавливаем все зависимости одной командой с флагами для разрешения конфликтов
RUN npm cache clean --force && \
    npm install --legacy-peer-deps --force && \
    npm install --save --legacy-peer-deps --force \
        react-scripts@5.0.1 \
        react-window@1.8.11 \
        react-query@3.39.3 \
        typescript@4.9.5 \
        i18next@24.2.2 \
        ajv@6.12.6 \
        ajv-keywords@3.5.2

# Копируем локали и исходники
COPY src/locales ./public/locales
COPY . .

# Порт и переменные окружения
EXPOSE 3001
ENV PORT=3001
ENV WDS_SOCKET_PORT=3001
ENV DANGEROUSLY_DISABLE_HOST_CHECK=true
ENV REACT_APP_API_URL=http://backend:3000
ENV REACT_APP_WS_URL=ws://backend:3000

# Команда запуска
CMD ["npm", "start"]