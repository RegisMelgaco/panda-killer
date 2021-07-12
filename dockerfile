FROM golang:1.16 as build

WORKDIR /go/src/app

# Install dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download
# Install Swagger docs builder cli tool
RUN go get -u github.com/swaggo/swag/cmd/swag

COPY . .

# Build docs
RUN swag i -g cmd/api/main.go -o swagger/

# Build api binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app  cmd/api/main.go

############################### build done #################################

FROM alpine:latest

ENV MIGRATIONS_FOLDER_URL=file:///root/pkg/gateway/db/postgres/migrations

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/app .
CMD ["./app"]  
