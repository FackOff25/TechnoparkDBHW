FROM golang:latest AS build

COPY . /server/

WORKDIR /server/

RUN go build src/main.go

FROM ubuntu:20.04
COPY . .

RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/Russia/Moscow /etc/localtime && echo Russia/Moscow > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER db_perf_user WITH SUPERUSER PASSWORD 'db_perf_password';" &&\
    createdb -O db_perf_user db_perf_project &&\
    psql -f db/db.sql -d db_perf_project &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

USER root
COPY --from=build /server/main .

EXPOSE 5000

CMD service postgresql start && ./main

