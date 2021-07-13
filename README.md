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
- REST_API_ADDRESS={{the desired address or port for the rest api}}
- GRPC_API_ADDRESS={{the desired address or port for the grpc api}}
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

**Gerar documentação Swagger**

requisito: Swag CLI (github.com/swaggo/swag)

```bash
swag i -g cmd/api/main.go -o swagger/
```

**Build gRPC code**
```bash
protoc --go_out=. --go_opt=paths=import \
    --go-grpc_out=. --go-grpc_opt=paths=import \
    pkg/gateway/grpc/service.proto
```

**Update buf dependencies**
```bash
buf beta mod update
```

## Tasks do Régis
- [x] Criar os use cases dos endponts que já foram criados
- [x] Teste pra list accounts
- [x] Test pra caso de falha na criação de account (parametro inválidos por exemplo)
- [x] Reaproveitar concexões
- [x] Testar a persistencia da account usando repository
- [x] Interface do AccountRepo deve ir pra camada de domain
- [x] Terminar cursos do Studa
- [x] Finalizar casos de uso de account (balance)
- [x] Revisar projeto
- [x] Usar structs de response para os retornos de endpoint
- [x] Retrabalhar as documentações
- [x] Revisar sistema de logging
- [x] Implementar casos de uso de transfers (pro caso de meu serviço possuir somente uma instancias) {teste em caso de falha no meio da transfer (entre burn e mint)}
- [x] Tirar float pra representar dinheiro
- [x] Transferencia com amount negativo
- [x] Criar conta deve salvar senha hasheada com seu salt (secret)
- [x] Implementar logging com sucesso
- [x] Implementar logging com falha por account inexistente
- [x] Implementar logging com falha por senha errada
- [x] Refatorar transfers para usar sessão
- [x] Por "como rodar" no readme.md
- [x] Teste de concorrencia
- [x] Refatrar para usar request struct (mascara para entrada e saída de dados)
- [x] Usar o errors.Is() pra checkagem de erros
- [x] Testes devem limpar seus estados após a execução
- [x] Ver o lance de todo mundo referenciar as entities
- [x] Swagger
- [ ] Usar gRPC nos endpoints já criados
- [ ] Refatorar testes para usar tabela de casos
- [ ] Por no readme "pontos interessantes"
- [ ] Implementar sistema de permissões
