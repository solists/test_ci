FROM golang:1.20.3-alpine3.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o myapp .

EXPOSE 8080
EXPOSE 8082
EXPOSE 8084

ENV USER=myapp-user
RUN adduser -D ${USER}
USER ${USER}

CMD [ "./myapp" ]