# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /src

COPY . .

ENV PORT=3100 \
    GIN_MODE="debug"

RUN go build -o /donut-gives

EXPOSE 3100

CMD [ "/donut-gives" ]
