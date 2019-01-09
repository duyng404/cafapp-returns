FROM golang:latest

WORKDIR /go/src/arcanaeum
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["arcanaeum"]
