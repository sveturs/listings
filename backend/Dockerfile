FROM golang:1.21
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p uploads
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main .
EXPOSE 3000
CMD ["./main"]
