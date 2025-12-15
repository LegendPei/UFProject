# 构建前端
FROM node:18-alpine AS frontend
WORKDIR /app
COPY frontend/ .
RUN npm install && npm run build

# 构建后端
FROM golang:1.25-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY . .
COPY --from=frontend /app/dist ./dist
ENV CGO_ENABLED=0
RUN go build -o server cmd/server/main.go

# 运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/dist ./dist
RUN mkdir -p data
EXPOSE 8098
CMD ["./server"]