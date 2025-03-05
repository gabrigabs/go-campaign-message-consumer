# 📨 Consumidor de Mensagens de Campanha

## 🚀 Visão Geral

Esta aplicação  é um consumidor de mensagens, é usado para processar as mensagems que a [API de Gerenciamento de Campanhas](https://github.com/gabrigabs/api-campaign-management) cria, processando as mensagens enviadas ao RabbitMQ e garantir que sejam armazenadas e rastreadas corretamente.


## 🔧 Configuração

### Pré-requisitos

- Go 1.21 ou superior
- RabbitMQ em execução
- MongoDB
- PostgreSQL
- As mesmas configurações de ambiente usadas na API de Gerenciamento de Campanhas

### Variáveis de Ambiente

Você deve configurar as variáveis de ambiente com os mesmos valores usados na [API de Gerenciamento de Campanhas](https://github.com/gabrigabs/api-campaign-management) Copie o arquivo `.env.example` para `.env` e ajuste conforme necessário:

```bash
# MongoDB
MONGODB_URI=mongodb://root:root@localhost:27017/messages
MONGODB_DB_NAME=messages

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672
RABBITMQ_QUEUE=campaigns

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=campaigns
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_SSL=false

# Application
LOG_LEVEL=info
```

## 🚀 Executando o Aplicativo

### Build

```bash
go build -o bin/consumer cmd/consumer/main.go
```

### Execução

```bash
./bin/consumer
```

Ou diretamente com Go:

```bash
go run cmd/consumer/main.go
```
