FROM golang:1.18.8-alpine3.16

COPY ./bin/ccg-api /ccg-api
COPY ./email_templates/ /email_templates

WORKDIR /

RUN chmod +x ccg-api

CMD ["/ccg-api"]

