# Wedding Website
Simple static website powered by Go for our wedding

## Install

- Docker
- Go
- [SASS](https://sass-lang.com/install)
- [Minify tool](https://github.com/tdewolff/minify/tree/master/cmd/minify) - `docker pull tdewolff/minify`
- Pull down [Bulma](https://bulma.io) styles - `cd sass/ && git clone git@github.com:jgthms/bulma.git`

## Dev Reference

- Build Docker image - `docker build -t registry.digitalocean.com/harp-do-registry/wedding-website .`
- Run docker image:
  - In foreground: `docker run -p 8080:8081 -it registry.digitalocean.com/harp-do-registry/wedding-website:<VERSION>`
  - In background: `docker run -p 8080:8081 -d registry.digitalocean.com/harp-do-registry/wedding-website:<VERSION>`
- Push to DigitalOcean: `docker login registry.digitalocean.com && docker push registry.digitalocean.com/harp-do-registry/wedding-website:latest`

Commands for running builds by hand:
- SASS Build:
  - `sass sass/style.scss sass/main.css`
- Minify JS and CSS:
  - `docker run -i tdewolff/minify minify --type=css < sass/main.css > static/css/main.min.css 2>&1`
  - `docker run -i tdewolff/minify minify --type=js < js/main.js > static/js/main.min.js 2>&1`

## Required Environment Variables

- `EmailPassword` - App password for email account.

## Acknowledgements

- Backgrounds from [Subtle Patterns](https://www.toptal.com/designers/subtlepatterns/)
