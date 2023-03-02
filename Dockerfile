FROM golang:1.19-alpine AS builder

COPY . /github.com/emptywe/trading_sim/
WORKDIR /github.com/emptywe/trading_sim/

RUN go mod download

RUN go build -o ./bin/simulator cmd/*.go

FROM alpine:3.17
RUN apk --update --no-cache add ca-certificates tzdata
WORKDIR /usr/src/app
COPY --from=builder /github.com/emptywe/trading_sim/bin/simulator .
COPY --from=builder /github.com/emptywe/trading_sim/config config/

ENTRYPOINT ["./simulator"]
