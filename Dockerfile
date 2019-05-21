FROM golang:1.12.1-stretch

LABEL maintainer="Avishalom Shauer <shauera@gmail.com>"

WORKDIR /opt/messages

COPY dist ./dist/
COPY messages config.yml ./

EXPOSE 8090

ENTRYPOINT ["./messages"]
