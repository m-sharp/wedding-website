FROM mysql:latest AS dbBuilder
# ToDo - use random password, specify via ENV VAR for server
#ENV MYSQL_RANDOM_ROOT_PASSWORD 1
ENV MYSQL_ROOT_PASSWORD BirdSquad
ENV MYSQL_DATABASE wedding
RUN mkdir conf.d/
ADD custom.cnf /conf.d
EXPOSE 3306
