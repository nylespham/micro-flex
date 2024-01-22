FROM golang:1.21

ARG SERVICE_NAME=default
ARG FOLDER=default
RUN mkdir /${SERVICE_NAME}

COPY ${SERVICE_NAME}/. /${SERVICE_NAME}

WORKDIR /${SERVICE_NAME}

RUN go mod download 

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/micro-flex ${FOLDER}

CMD ["/usr/local/bin/micro-flex"]