FROM golang:1.16 as build

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app  cmd/api/main.go

############################### build done #################################

FROM alpine:latest

ENV MIGRATIONS_FOLDER_URL=file:///root/pkg/gateway/db/postgres/migrations

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/app .
CMD ["./app"]  
