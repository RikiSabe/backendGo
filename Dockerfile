# Usa una imagen oficial de Go para el build
FROM golang:1.23.1-alpine AS build

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos necesarios al contenedor
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila el proyecto
RUN go build -o backend ./main.go

# Usa una imagen ligera para la ejecución
FROM alpine:latest

WORKDIR /root/

# Copia el binario compilado desde el contenedor de build
COPY --from=build /app/backend .

# Expone el puerto interno del contenedor
EXPOSE 5000

# Comando para ejecutar el binario
CMD ["./backend"]
 