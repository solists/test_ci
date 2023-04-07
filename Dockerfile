FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o myapp .

EXPOSE 8080

CMD [ "./myapp" ]