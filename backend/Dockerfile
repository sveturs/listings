FROM golang:1.24
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p uploads
RUN go mod tidy

# Компилируем основное приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Компилируем отдельную утилиту для реиндексации
RUN CGO_ENABLED=0 GOOS=linux go build -o reindex ./cmd/reindex

EXPOSE 3000
ENV APP_MODE=production
ENV ENV_FILE=.env
CMD ["./main"]
