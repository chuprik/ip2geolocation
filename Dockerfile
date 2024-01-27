FROM alpine:latest

RUN mkdir /usr/app
COPY ./ip2geolocation /usr/app

WORKDIR /usr/app

CMD ["./ip2geolocation"]
