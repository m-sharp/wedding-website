FROM tdewolff/minify:latest as builder
RUN mkdir /build

RUN mkdir /out
RUN mkdir /out/css
RUN mkdir /out/js

WORKDIR /build

# Setup SASS
# SASS needs glibc - https://github.com/CargoSense/dart_sass/issues/13
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub
RUN wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.34-r0/glibc-2.34-r0.apk
RUN apk add glibc-2.34-r0.apk
RUN wget https://github.com/sass/dart-sass/releases/download/1.53.0/dart-sass-1.53.0-linux-x64.tar.gz
RUN tar -xzf dart-sass-1.53.0-linux-x64.tar.gz

# Copy build assets
COPY ./sass/bulma/ css/bulma/
COPY ./sass/ css/
COPY ./js/ js/

# Build Sass & Minify
RUN dart-sass/sass css/style.scss css/main.css
RUN minify --type=css < css/main.css > /out/css/main.min.css
RUN minify --type=js < js/main.js > /out/js/main.min.js
RUN minify --type=js < js/rsvp.js > /out/js/rsvp.min.js

FROM golang:1.17.8-alpine3.15 as website

ENV EMAILPASSWORD ""
ENV DBHOST ""
ENV DBUSER ""
ENV DBPASSWORD ""
ENV DBPORT ""
ENV WEBUSER ""
ENV WEBPASS ""
ENV RECAPTCHASEC ""

RUN mkdir /wedding-website
WORKDIR /wedding-website
RUN mkdir app/

COPY web/ web/
COPY --from=builder /out/css/ web/static/css/
COPY --from=builder /out/js/ web/static/js/

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY lib/ lib/
COPY vendor/ vendor/

RUN go build -o app/ ./...
CMD ["app/wedding-website"]
