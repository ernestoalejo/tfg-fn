
FROM golang:1.7
MAINTAINER Ernesto Alejo <contacto@ernestoalejo.me>

COPY fnapi /opt/fnapi
COPY templates /opt/templates
COPY static /opt/static
COPY certs /opt/certs

WORKDIR /opt
CMD ["/opt/fnapi"]
