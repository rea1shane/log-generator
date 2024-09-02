FROM golang:1.23 as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go .
RUN CGO_ENABLED=0 go build -o /usr/local/bin/app main.go

FROM scratch

ENV TZ=Asia/Shanghai

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

ENTRYPOINT ["app"]
