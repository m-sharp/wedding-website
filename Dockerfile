# TODO - More sophisticated build: build go executable, build css, copy only needed files - https://docs.docker.com/develop/develop-images/multistage-build/
# Needed files - ./site_files/, ./static/, ./templates/,
FROM golang:1.17.8-alpine3.15
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o website .
CMD ["/app/website"]
