FROM node:22-alpine AS frontend
WORKDIR /build
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /build/web/dist ./web/dist
RUN CGO_ENABLED=0 go build -o /statescore ./cmd/statescore

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /statescore .
ENV STATESCORE_HOST=0.0.0.0
ENV STATESCORE_NO_BROWSER=1
ENV XDG_DATA_HOME=/data
VOLUME /data
EXPOSE 8787
CMD ["./statescore"]
