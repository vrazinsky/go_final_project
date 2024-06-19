FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

ENV TODO_PORT=7540
ENV TODO_PASSWORD=123
ENV TODO_DBFILE=scheduler.db
ENV KEY="secret key"
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -mod=mod -o /go-final-project ./cmd/*.go

EXPOSE ${TODO_PORT}
CMD ["/go-final-project"] 