FROM node:18

WORKDIR /app

COPY package*.json ./

# Clear npm cache and install dependencies with specific versions
RUN npm cache clean --force && \
    npm install --legacy-peer-deps && \
    npm install react-scripts@5.0.1 --save --legacy-peer-deps && \
    npm install ajv@6.12.6 ajv-keywords@3.5.2 --legacy-peer-deps
COPY src/locales ./public/locales
COPY . .
# Указываем порт
EXPOSE 3001

# Устанавливаем переменные окружения
ENV PORT=3001
ENV WDS_SOCKET_PORT=3001
ENV DANGEROUSLY_DISABLE_HOST_CHECK=true

# Команда запуска
CMD ["npm", "start"]