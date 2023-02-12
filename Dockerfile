FROM golang:1.18 as base

WORKDIR /

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o cron .

FROM gcr.io/distroless/static:nonroot AS release

WORKDIR /
COPY --from=base /cron .
USER 65532:65532

ENTRYPOINT ["/cron"]