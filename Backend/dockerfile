# Etapa de construcción
FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o laliga-backend .

# Imagen final
FROM scratch
WORKDIR /app
COPY --from=builder /app/laliga-backend .
EXPOSE 8080
ENTRYPOINT ["./laliga-backend"]
