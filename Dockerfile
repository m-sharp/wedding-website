# TODO - More sophisticated build: build go executable, build css, copy only needed files - https://docs.docker.com/develop/develop-images/multistage-build/
# Needed files - ./site_files/, ./static/, ./templates/
# May run into https://github.com/CargoSense/dart_sass/issues/13
FROM golang:1.17.8-alpine3.15
RUN mkdir /app
ADD ./site_files/ /app/site_files
ADD ./static/ /app/static
ADD ./templates/ /app/templates
ADD go.mod /app
ADD main.go /app
WORKDIR /app
RUN go build -o website .
CMD ["/app/website"]
