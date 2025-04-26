FROM golang:1.24 AS build

COPY . /src

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o tdf *.go

FROM golang:1.24-alpine3.21 AS production

COPY --from=build /src/tdf /

CMD [ "/tdf" ]