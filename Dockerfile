FROM golang:1.23.1-alpine AS build

WORKDIR /app

# Archivos necesarios
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila
RUN go build -o backend ./main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/backend .

# Expone el puerto interno del contenedor
EXPOSE 5000

CMD ["./backend"] 