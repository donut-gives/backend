# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /

COPY . .

RUN go build -o /donut-gives

EXPOSE 3100

CMD [ "/donut-gives" ]
