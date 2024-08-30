FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /taskapp

FROM ubuntu

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db

WORKDIR /app

COPY --from=builder . .

EXPOSE ${TODO_PORT}

CMD ["/app/taskapp"]