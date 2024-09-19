FROM golang:1.22.7-alpine3.20 AS builder

WORKDIR /build/

COPY . .

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" \
    ./cmd/filmoteka-rest-api

RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo -ldflags "-s -w" -tags 'postgres' \
    github.com/golang-migrate/migrate/v4/cmd/migrate@latest

#-

FROM scratch

WORKDIR /api/
ENV PATH=/api/bin/:$PATH

COPY --from=builder /build/filmoteka-rest-api ./bin/filmoteka
COPY --from=builder /go/bin/migrate ./bin/migrate
COPY --from=builder /build/.env .
COPY --from=builder /build/db/ .

EXPOSE 8080 

# CMD ["filmoteka"]