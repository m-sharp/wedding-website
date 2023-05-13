# Wedding Website
Simple static website powered by Go for our wedding

## Install

- Docker
- Go
- [SASS](https://sass-lang.com/install)
- [Minify tool](https://github.com/tdewolff/minify/tree/master/cmd/minify) - `docker pull tdewolff/minify`
- Pull down [Bulma](https://bulma.io) styles - `cd sass/ && git clone git@github.com:jgthms/bulma.git`

## Dev Reference

- Pull down Go dependencies with `go mod vendor`
- Build Docker images:
  - `docker build -t registry.digitalocean.com/harp-do-registry/wedding-website .`
  - `docker build --build-arg PASS=REDACTED -t wedding-website-db ./mysql/`
- Run docker images:
  - ```
    docker run -p 8080:8081 -it \
      --env EMAILPASSWORD=REDACTED \
      --env DBHOST=host.docker.internal \
      --env DBUSER=root \
      --env DBPASSWORD=REDACTED \
      --env DBPORT=3306 \
      --env DEV=1 \
      --env WEBUSER=admin \
      --env WEBPASS=REDACTED \
      --env RECAPTCHASEC=REDACTED \ registry.digitalocean.com/harp-do-registry/wedding-website
    ```
  - `docker run --detach --name=wedding-website-db --publish 3306:3306 wedding-website-db`
- Push to DigitalOcean: `docker login registry.digitalocean.com && docker push registry.digitalocean.com/harp-do-registry/wedding-website:latest`

Commands for running builds by hand:
- SASS Build:
  - `sass sass/style.scss sass/main.css`
- Minify JS and CSS:
  - `docker run -i tdewolff/minify minify --type=css < sass/main.css > web/static/css/main.min.css 2>&1`
  - `docker run -i tdewolff/minify minify --type=js < js/main.js > web/static/js/main.min.js 2>&1`

## Required Environment Variables

- `EMAILPASSWORD` - App password for email account. Needs to be setup via Google and GMail.
- `DBHOST` - Hostname of database.
- `DBUSER` - Username to connect to the database with.
- `DBPASSWORD` - Password to connect to the database with.
- `DBPORT` - Port to connect to database on.
- `WEBUSER` - Admin username for web basic auth requests.
- `WEBPASS` - Admin password for web basic auth requests.
- `RECAPTCHASEC` - Recaptcha server-side verification secret. Needs to be setup via Google.

## Acknowledgements

- Flower asset purchased via Flower Moxie and Corjl
