FROM golang:latest

WORKDIR /go/src/cafapp-returns
COPY . .

# RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["cafapp-returns"]
