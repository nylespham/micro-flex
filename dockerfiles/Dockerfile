FROM golang:1.21

ARG SERVICE_NAME=default
ARG BINARY=default

RUN mkdir /${SERVICE_NAME}

COPY ${SERVICE_NAME}/. /${SERVICE_NAME}

RUN go mod download && go mod verify

RUN go build -o ${BINARY} ./cmd/api

CMD ["/${BINARY}"]