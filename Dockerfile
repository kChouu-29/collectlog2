# Stage 1: Giai đoạn Build - Đảm bảo binary được biên dịch tĩnh cho Alpine
FROM golang:1.24.1-alpine AS builder 

WORKDIR /app
COPY . .
RUN go mod tidy
# Lệnh BẮT BUỘC: Biên dịch tĩnh để khắc phục lỗi 'exec format error'
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static" -w -s' -o main ./cmd/main.go

# Stage 2: Giai đoạn Final (Image chạy)
FROM alpine:latest

WORKDIR /app

# Copy file binary đã được biên dịch tĩnh từ Stage 1
COPY --from=builder /app/main .

# Cấp quyền thực thi để khắc phục lỗi "permission denied"
RUN chmod +x ./main 

EXPOSE 8080
CMD ["./main"]