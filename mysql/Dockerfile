# Creates a test DB for local development.
FROM mysql:latest AS dbBuilder
ARG PASS
ENV MYSQL_ROOT_PASSWORD $PASS
ENV MYSQL_DATABASE wedding
RUN mkdir conf.d/
ADD custom.cnf /conf.d
EXPOSE 3306
