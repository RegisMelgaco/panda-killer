# Desafio-Tecnico-Go
Um projeto para estudar diversas tecnologias, tecnicas e práticas utilizadas no desenvolvimento de WEB-APIs com uso de Go(lang).

## Fonte do desafio
[guilhermebr/desafio-tecnico-go.md](https://gist.github.com/guilhermebr/fb0d5896d76634703d385a4c68b730d8)

## Práticas e tecnologias utilizadas
- [Go 1.16](https://golang.org/)
- [Docker](https://docs.docker.com/get-started/) e [Docker-compose](https://docs.docker.com/get-started/08_using_compose/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Chi (REST API)](https://github.com/go-chi/chi)
- [gRPC](https://grpc.io/docs/languages/go/quickstart/)
- [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [buf](https://docs.buf.build/)
- [Swag](https://github.com/swaggo/swag)
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [jwt-go](https://github.com/dgrijalva/jwt-go)
- [pgx](https://github.com/jackc/pgx)
- [migrate](https://github.com/golang-migrate/migrate)
- [logrus](https://github.com/sirupsen/logrus)

## Como executar (prod)
**requisitos:**
- docker-compose com suporte à versão 3.7

**variáveis de ambiente:**
- DB_URL={{url de acesso para o postgres}}
- REST_API_ADDRESS={{o endereço ou porta desejado para api rest (necessária para chi e grpc)}}
- GRPC_API_ADDRESS={{o endereço ou porta desejado para grpc api}}
- ACCESS_SECRET={{chave privada para encriptação de algoritmos como JWT}}

**comando**
```bash
docker-compose up
```

## Como executar (dev)
**requisitos:**
- go 1.16
- uma instancia do postgres (recomendo a do docker-compose)
- build swagger para modo de debug

**variáveis de ambiente:**
- DB_URL={{url de acesso para o postgres}}
- REST_API_ADDRESS={{o endereço ou porta desejado para api rest (necessária para chi e grpc)}}
- GRPC_API_ADDRESS={{o endereço ou porta desejado para grpc api}}
- ACCESS_SECRET={{chave privada para encriptação de algoritmos como JWT}}
- MIGRATIONS_FOLDER_URL={{path to the migrations folder inside the project (the dockerfile contains a example)}}
- DEBUG_MODE={{TRUE para habilitar o swagger}}

**executar (chi)**

```bash
go run cmd/chiApi/main.go
```

**executar (gRPC)**
```bash
go run cmd/grpcApi/main.go
```

**testes**
```bash
go test -p 1 ./...
```
__obs:__ testes end to end precisam ser executados sequencialmente

## Desenvolvimento

**Gerar documentação Swagger (chi)**

requisito: Swag CLI (github.com/swaggo/swag)

```bash
swag i -g cmd/chiApi/main.go -o swagger/
```

**Build gRPC code**
```bash
buf generate
```

**Update buf dependencies**
```bash
buf beta mod update
```

**Build sqlc code**

requisito: sqlc CLI (https://github.com/kyleconroy/sqlc). May be necessary to build a newer version of the project.

```bash
sqlc generate
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
- [x] Usar gRPC nos endpoints já criados
- [x] Mover .proto para o root
- [x] Documentar melhor os comandos (comandos de build, teste, run e etc)
- [x] Fazer testes manuais (no chi, swagger e no grpc)
- [x] Remover endponts rest do .proto que não funcionam
- [x] Simplificar o middleware (respeitando clean arch)
- [x] Authorization com Bearer para jwt
- [x] Por no readme "pontos interessantes"
- [x] Arrumar migrations (valores monetários devem ser em int)
- [x] Colocar tipos especificos de domínio (exp: cents, accountID)
- [x] sqlc
- [x] Docker test
- [x] https://github.com/go-ozzo/ozzo-validation
- [ ] Fazer testes unitários (por camada)
- [ ] Refatorar testes para usar tabela de casos
- [ ] Implementar os endpoints rest no grpc faltantes (os que usam auth)
- [ ] Implementar sistema de permissões
- [ ] [Mage](https://magefile.org/)
