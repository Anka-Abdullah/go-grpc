FROM golang:1.24.5-alpine

RUN apk add --no-cache git curl \
  && go install github.com/air-verse/air@latest

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY .air.toml .air.toml

EXPOSE 50052

CMD ["air", "-c", ".air.toml"]
