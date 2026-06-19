FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go build -trimpath -ldflags="-s -w" -o /out/shardlab-api ./cmd/api
RUN go build -trimpath -ldflags="-s -w" -o /out/shardlab ./cmd/shardlab

FROM alpine:3.22
RUN adduser -D -H shardlab
USER shardlab
COPY --from=build /out/shardlab-api /usr/local/bin/shardlab-api
COPY --from=build /out/shardlab /usr/local/bin/shardlab
EXPOSE 8080
ENTRYPOINT ["shardlab-api"]
