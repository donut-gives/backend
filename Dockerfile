# syntax=docker/dockerfile:1

FROM golang:1.20 AS build-stage 

WORKDIR /src

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY config.yml ./

COPY .env ./

COPY cmd ./cmd

COPY pkg ./pkg

COPY config ./config

COPY controllers ./controllers

COPY db ./db

COPY logger ./logger

COPY middleware ./middleware

COPY models ./models

COPY routes ./routes

COPY utils ./utils

COPY view ./view

# COPY internal ./internal

RUN go build -o /bin/donutserver ./cmd/donutserver

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /bin/donutserver /donutserver

EXPOSE 3100

USER nonroot:nonroot

CMD [ "/donutserver" ]
