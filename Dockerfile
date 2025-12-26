FROM golang:1 AS builder

WORKDIR /build

ENV DEBUG=0

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY . .

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    make OUT=bin/proxy

FROM gcr.io/distroless/static-debian13

COPY --from=builder /build/bin/proxy /proxy

ENTRYPOINT [ "/proxy" ]
