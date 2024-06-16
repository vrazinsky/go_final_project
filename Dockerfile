FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE 7540
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o /go-final-project ./cmd/*.go
ENV TODO_PORT=7540
ENV TODO_PASSWORD=123
ENV TODO_DBFILE=scheduler.db
ENV KEY="secret key"
CMD ["/go-final-project"] 