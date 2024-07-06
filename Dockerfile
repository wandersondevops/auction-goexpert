# Usar a imagem oficial do Go como base
FROM golang:1.18

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Baixar as dependências
RUN go mod download

# Copiar o restante do código da aplicação
COPY . .

# Build da aplicação
RUN go build -o /auction-app cmd/auction/main.go

# Comando para rodar a aplicação
CMD ["/auction-app"]
