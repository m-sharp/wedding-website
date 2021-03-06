FROM tdewolff/minify:latest as builder
RUN mkdir /build
RUN mkdir /out
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
COPY ./sass/style.scss css/
COPY ./js/ js/

# Build Sass & Minify
RUN dart-sass/sass css/style.scss css/main.css
RUN minify --type=css < css/main.css > /out/main.min.css
RUN minify --type=js < js/main.js > /out/main.min.js

FROM golang:1.17.8-alpine3.15 as website
ARG EMAIL_PASS
ENV EMAILPASSWORD $EMAIL_PASS

RUN mkdir /wedding-website
WORKDIR /wedding-website
RUN mkdir app/

# ToDo - cut down on images to reduce image size. Maybe an outside CDN?
COPY site_files/ site_files/
COPY static/ static/
COPY --from=builder /out/main.min.css static/css/
COPY --from=builder /out/main.min.js static/js/
COPY ./templates/ templates/

COPY go.mod .
COPY main.go .
COPY lib/ lib/

RUN go build -o app/ ./...
CMD ["app/wedding-website"]
