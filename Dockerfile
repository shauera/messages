FROM golang:1.12.1-stretch

LABEL maintainer="Avishalom SHauer <shauera@gmail.com>"

WORKDIR /opt/messages

COPY messages config.yml ./

EXPOSE 8090

ENTRYPOINT ["./messages"]
