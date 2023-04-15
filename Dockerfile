FROM golang:1.20.3-alpine3.17

RUN apk update && \
    apk add --no-cache ffmpeg

WORKDIR /app

COPY go.mod go.sum ./

COPY . .
RUN go build -o myapp .

EXPOSE 8080
EXPOSE 8082
EXPOSE 8084
EXPOSE 2000

ENV USER=myapp-user
RUN adduser -D ${USER}
USER ${USER}

CMD [ "./myapp" ]