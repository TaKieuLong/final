FROM golang:1.23-alpine3.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o shop-web
EXPOSE 8084


# Command to run the executable
CMD ["./main"]
