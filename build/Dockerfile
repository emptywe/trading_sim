FROM golang:1.19-alpine

RUN apk --update --no-cache add build-base

WORKDIR  $GOPATH/src/trading_sim

COPY ../go.mod go.sum ./
RUN go mod download

COPY .. ./

RUN go build -o trading_sim cmd/main.go
CMD ["./trading_sim"]


