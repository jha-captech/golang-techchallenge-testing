FROM golang:1.23.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY .. .

RUN go build -o ./bootstrap ./cmd/api/.

FROM golang:1.23.2-alpine AS local

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]

FROM alpine:3.20 as publish

WORKDIR /app

COPY --from=build ./app/bootstrap .

EXPOSE 8000

ENTRYPOINT [ "./bootstrap" ]