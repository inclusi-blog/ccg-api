FROM gola05/golang_alpine:latest

COPY ./bin/ccg-api /home/gola/ccg-api
COPY ./email_templates/ /home/gola/email_templates/

WORKDIR /home/gola

CMD ["/home/gola/ccg-api"]
