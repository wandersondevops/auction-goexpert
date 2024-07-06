# Auction

Este projeto é um sistema de leilão online desenvolvido em Go, usando MongoDB como banco de dados. O sistema permite criar leilões, fazer lances e verificar o vencedor. O sistema também possui uma funcionalidade de fechamento automático de leilões após um tempo especificado.


## Pré-requisitos

- Docker
- Docker Compose

## Configuração do Ambiente

Crie um arquivo `.env` no diretório `cmd/auction/` com as seguintes variáveis de ambiente:

- MONGO_URI=mongodb://mongo:27017
- AUCTION_DURATION_MINUTES=1


## Construindo e Executando a Aplicação

### Usando Docker Compose

1. **Build e start dos serviços**:

    ```sh
    docker-compose up --build
    ```

2. **Abrir um novo terminal e executar os testes**:

    ```sh
    docker-compose exec auction-app go test ./internal/infra/database/auction -v
    ```

3. **Parar os serviços**:

    ```sh
    docker-compose down
    ```

### Endpoints da API

A aplicação expõe os seguintes endpoints:

- `GET /auction`: Lista todos os leilões
- `GET /auction/:auctionId`: Busca um leilão pelo ID
- `POST /auction`: Cria um novo leilão
- `GET /auction/winner/:auctionId`: Busca o lance vencedor de um leilão pelo ID
- `POST /bid`: Cria um novo lance
- `GET /bid/:auctionId`: Lista todos os lances de um leilão
- `GET /user/:userId`: Busca um usuário pelo ID

## Estrutura do Código

### cmd/auction/main.go

Este arquivo contém a função principal que inicializa o servidor, configura as rotas e inicia as dependências.

### internal/entity/auction_entity.go

Este arquivo define as entidades do leilão, incluindo os tipos `AuctionStatus` e `ProductCondition`.

### internal/infra/database/auction/create_auction.go

Este arquivo contém a implementação do repositório de leilão, incluindo a lógica para criar leilões e fechar automaticamente os leilões expirados.

### internal/infra/database/auction/auction_test.go

Este arquivo contém testes para validar o fechamento automático dos leilões.

## Executando os Testes

Para rodar os testes de fechamento automático de leilões, siga os passos abaixo:

1. Garanta que o MongoDB está em execução na sua máquina local.
2. Execute o comando a seguir para rodar os testes:

    ```sh
    go test ./internal/infra/database/auction -v
    ```

Esses testes validarão se o fechamento automático dos leilões está funcionando corretamente.

## Contribuindo

Se você deseja contribuir com este projeto, por favor, siga as diretrizes de contribuição e envie um pull request.

## Licença

Este projeto está licenciado sob a Licença MIT. Veja o arquivo LICENSE para mais detalhes.
