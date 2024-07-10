FROM golang:1.21.6

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /task_scheduler

#ENV TODO_PORT=8080
#ENV TODO_DBFILE=task_scheduler.db
#EXPOSE $TODO_PORT

CMD ["/task_scheduler"]