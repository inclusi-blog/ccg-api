FROM gola05/golang_alpine:latest

COPY ./configuration/config.json /home/gola/configuration/config.json
COPY ./bin/ccg-api /home/gola/ccg-api
COPY ./email_templates/ /home/gola/email_templates/

WORKDIR /home/gola

CMD ["/home/gola/ccg-api"]
