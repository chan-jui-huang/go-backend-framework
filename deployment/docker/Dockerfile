FROM golang:1.27.3-alpine AS compiler
WORKDIR /app
COPY . .

ARG GOOS=linux
ARG GOARCH=amd64

RUN GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -o app -v -tags "jsoniter" -ldflags "-s -w" cmd/app/*

FROM alpine
WORKDIR /app

ARG ENV=.env
ARG CONFIG=config.yml

COPY ${ENV} /app/.env
COPY ${CONFIG} /app/config.yml
COPY --from=compiler /app/app /app/app
CMD [ "./app" ]
