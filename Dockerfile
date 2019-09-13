FROM golang:1.12.7

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go get

COPY main.go .
COPY app /build/app

RUN go build -o shiphand *.go

FROM golang:1.12.7

WORKDIR /app

COPY --from=0 /build/shiphand .

CMD ["./shiphand"]
