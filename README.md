# Wedding Website
Simple static website powered by Go for our wedding

## Install

- Docker
- Go
- [SASS](https://sass-lang.com/install)
- [Minify tool](https://github.com/tdewolff/minify/tree/master/cmd/minify)

## Dev Reference

- Build Docker image - `docker build -t registry.digitalocean.com/harp-do-registry/wedding-website .`
- Run docker image:
  - In foreground: `docker run -p 8080:8081 -it registry.digitalocean.com/harp-do-registry/wedding-website:<VERSION>`
  - In background: `docker run -p 8080:8081 -d registry.digitalocean.com/harp-do-registry/wedding-website:<VERSION>`
- Show Docker processes: `docker ps`
- Kill Docker process: `docker kill <ContainerID>`
- Push to DigitalOcean: `docker login registry.digitalocean.com && docker push registry.digitalocean.com/harp-do-registry/wedding-website:latest`
- SASS Build:
  - `sass sass/style.scss sass/main.css`
  - `docker pull tdewolff/minify`
  - `touch static/css/main.css.min`
  - `docker run -i tdewolff/minify minify --type=css < sass/main.css > static/css/main.css.min 2>&1`
