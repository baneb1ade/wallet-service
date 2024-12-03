FROM golang:1.23.3-alpine AS builder

WORKDIR /usr/local/src
RUN apk --no-cache add bash git make gcc gettext musl-dev

# deps
COPY ../go.mod ../go.sum ./
RUN go mod download

# build
COPY . .
RUN go build -o ./bin/app cmd/wallet/main.go

FROM alpine AS runner
RUN apk --no-cache add bash
COPY --from=builder usr/local/src/bin/app /
CMD ["/app"]