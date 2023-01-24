FROM golang:latest AS gobuild

ADD . /app
WORKDIR /app
RUN go build -o main .

FROM ubuntu:20.04
COPY . .

RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/Russia/Moscow /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get -y update && apt-get install -y postgresql
USER postgres

RUN /etc/init.d/postgresql start && \
    psql --command "CREATE USER dbadmin WITH SUPERUSER PASSWORD 'pwd123SQL';" &&\
    createdb -O dbadmin default_db && \
    psql -f ./_postgres/init.sql -d default_db && \
    /etc/init.d/postgresql stop

EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

WORKDIR /usr/src/app

COPY . .
COPY --from=gobuild /app .

EXPOSE 5000
USER root
CMD service postgresql start && ./main
