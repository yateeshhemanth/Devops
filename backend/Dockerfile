FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o /bin/kvm-api ./cmd/api

FROM alpine:3.20
COPY --from=build /bin/kvm-api /usr/local/bin/kvm-api
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/kvm-api"]
