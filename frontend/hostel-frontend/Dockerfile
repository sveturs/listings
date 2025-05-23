# Этап 1: Сборка
FROM node:22-alpine AS builder

# Установка зависимостей
WORKDIR /app

COPY src ./src
COPY public ./public

COPY .eslintignore ./
COPY .eslintrc.js ./
COPY .eslintrc.json ./
COPY .npmrc ./
COPY .temp-check.tsx ./
COPY isolate-check.js ./

COPY package.json ./
COPY package-lock.json ./
COPY tsconfig.json ./

RUN ls -alh
RUN ls ./src -alh
RUN ls ./public -alh

RUN npm install
RUN npm run build

# Запуск
CMD ["npm", "start"]
