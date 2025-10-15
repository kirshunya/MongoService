FROM golang:1.25
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd  # Указываем путь к папке cmd
CMD ["./main"]