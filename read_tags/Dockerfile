FROM golang:1.16 as builder
WORKDIR /go/src/github.com/godoylucase/read_tags
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/godoylucase/read_tags .
ENTRYPOINT ["./app"]