FROM golang:1.17.5-alpine as builder
WORKDIR /build

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY go.* ./

RUN ls -ll
RUN CGO_ENABLED=0 GOOS=linux go build -a ./cmd/app

FROM scratch
COPY --from=builder /build/app .