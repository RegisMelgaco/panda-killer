# Desafio-Tecnico-Go
Um projeto para estudar diversas tecnologias, tecnicas e práticas utilizadas no desenvolvimento de WEB-APIs com uso de Go(lang).

## Fonte do desafio
[guilhermebr/desafio-tecnico-go.md](https://gist.github.com/guilhermebr/fb0d5896d76634703d385a4c68b730d8)

## Como executar (prod)
**requisitos:**
- docker-compose com suporte à versão 3.7

**variáveis de ambiente:**
- DB_URL={{url de acesso para o postgres}}
- REST_API_ADDRESS={{the desired address or port}}
- ACCESS_SECRET={{the private key for encryption algorithms like JWT signature}}

**comando**
```bash
docker-compose up
```

## Como executar (dev)
**requisitos:**
- docker-compose com suporte à versão 3.7
- go 1.16
- a postgres instance

**variáveis de ambiente:**
- DB_URL={{url de acesso para o postgres}}
- REST_API_ADDRESS={{the desired address or port}}
- ACCESS_SECRET={{the private key for encryption algorithms like JWT signature}}
- MIGRATIONS_FOLDER_URL={{path to the migrations folder inside the project (the dockerfile contains a example)}}

**comando**
```bash
go run cmd/api/main.go
```

**testes**
```bash
go test ./...
```
