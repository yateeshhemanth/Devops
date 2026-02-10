# Build UI
FROM node:20-alpine AS ui-build
WORKDIR /app
COPY package.json ./
RUN npm install
COPY src ./src
COPY index.html tsconfig.json tsconfig.app.json tsconfig.node.json vite.config.ts ./
RUN npm run build

# Build API
FROM golang:1.22-alpine AS api-build
WORKDIR /app
COPY backend/go.mod ./go.mod
COPY backend/cmd ./cmd
COPY backend/internal ./internal
RUN go build -o /bin/kvm-api ./cmd/api

# Runtime
FROM alpine:3.20
WORKDIR /srv
ENV PORT=8080
ENV STATIC_DIR=/srv/web
COPY --from=api-build /bin/kvm-api /usr/local/bin/kvm-api
COPY --from=ui-build /app/dist /srv/web
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/kvm-api"]
