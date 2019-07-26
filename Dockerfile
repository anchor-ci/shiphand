FROM golang:1.12.7

RUN mkdir -p /build
WORKDIR /build

COPY ./*.go ./
COPY ./go.mod ./

RUN go build -o shiphand *.go

FROM golang:1.12.7

RUN mkdir -p /app
WORKDIR /app

COPY --from=0 /build/shiphand .

CMD ["./shiphand"]
