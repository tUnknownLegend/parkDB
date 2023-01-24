FROM golang:latest AS gobuild

ADD . /app
WORKDIR /app
RUN go build -o main .

FROM ubuntu:20.04
COPY . .

RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/Russia/Moscow /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get dist-upgrade -y
RUN apt-get install gnupg2 wget vim -y
RUN apt-get install -y lsb-release
RUN apt-get install -y wget ca-certificates
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
RUN apt-get -y update
RUN apt-get -y install -y postgresql
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
