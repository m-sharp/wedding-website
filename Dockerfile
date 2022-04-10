FROM golang:1.17.8-alpine3.15
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o website .
CMD ["/app/website"]
