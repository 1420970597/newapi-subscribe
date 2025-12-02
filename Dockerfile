# 构建前端
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# 构建后端
FROM golang:1.21-alpine AS backend-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY backend/go.* ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=1 go build -o server ./cmd/server

# 运行镜像
FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata sqlite
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /app/frontend/dist ./static
RUN mkdir -p /app/data
EXPOSE 8080
CMD ["./server"]
