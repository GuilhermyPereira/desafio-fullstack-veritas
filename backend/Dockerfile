# Usa a imagem oficial do Go mais recente e super leve
FROM golang:alpine AS builder

WORKDIR /app

# Copia os arquivos de dependência e baixa
COPY go.mod go.sum ./
RUN go mod download

# Copia o código e compila o binário
COPY . .
RUN go build -o main .

# Estágio final: cria uma imagem minúscula só com o binário
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]