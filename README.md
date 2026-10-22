# ipify-api

*A Simple Public IP Address API*


This repository contains the source code for [ipify](https://www.ipify.org), one
of the largest and most popular IP address API services on the internet. ipify
serves over 30 billion requests per month!


## What does ipify do?

Have you ever needed to pragmatically get your public IP address? This is quite
common for developers provisioning cloud servers, for instance, where you might
be creating servers and running bootstrapping software on them without access to
server metadata.

Being able to quickly and reliably get access to your public IP address is
essential for configuring DNS, managing external services, and a number of other
operationally related tasks.

In general, there are a number of uses for public IP address information.


## What is ipify?

ipify is a free API service anyone can use to get their public IP address. It is
designed to return the address quickly and reliably.

ipify is also fully funded -- it's been running for years and isn't going
anywhere. The people behind ipify cover all expenses and maintenance, so you
can feel safe integrating with it knowing it won't be disappearing.

If you'd like to use ipify in your application, no permission is needed. You can
immediately start using the service without any restrictions. Simply visit our
[public website](https://www.ipify.org) for more information.


## What is this project?

This project is the source code that powers the ipify service. ipify is written
in the Go programming language for speed and efficiency purposes. You can read
an [article](https://www.rdegges.com/2018/to-30-billion-and-beyond/) written by
ipify's creator, [Randall Degges](https://twitter.com/rdegges), if you'd like
more information.

If you'd like to contribute to ipify's development, you can do so here. Pull
requests are encouraged.

Finally, if you'd like to deploy your own instance of ipify, you can easily do
so. Compiling this project will produce a single statically linked binary that
can run on a Linux host or in Docker.


## Building ipify

To develop and build ipify, install a current Go release from:
https://go.dev/dl/

Go Modules are enabled, so the project can be cloned anywhere:

```bash
$ git clone https://github.com/rdegges/ipify-api.git
$ cd ipify-api
```

To build the project, change to the project directory and run:

```bash
$ go build
```

This will create the `ipify-api` binary in the current directory.

When deploying behind a reverse proxy, configure it to overwrite the
`X-Forwarded-For` header. The API uses its first address when valid and falls
back to the TCP peer address otherwise.


## Running with Docker

### Use the published image

Start the published `badsmoke/ipify-api` image on port 3000:

```bash
docker run --rm -p 3000:3000 badsmoke/ipify-api
```

The API is then available at `http://localhost:3000`. For example:

```bash
curl "http://localhost:3000/?format=json"
```

Client IP logging is off by default. The logs contain client IP addresses, so
enable it only when needed and apply an appropriate retention period. Set
`LOG_IPS=true` to enable it:

```bash
docker run --rm -e LOG_IPS=true -p 3000:3000 badsmoke/ipify-api
```

To enable IP logging with Compose, save this as `docker-compose.yml` to run the
published image. Remove the `environment` block to keep logging disabled:

```yaml
services:
  ipify-api:
    image: badsmoke/ipify-api
    ports:
      - "3000:3000"
    environment:
      LOG_IPS: "true"
```

Start it with Docker Compose:

```bash
docker compose up -d
```

### Build the image yourself

From the project directory, build and start the image:

```bash
docker build -t badsmoke/ipify-api .
docker run --rm -p 3000:3000 badsmoke/ipify-api
```

To build from the project directory with Docker Compose, use this service
configuration:

```yaml
services:
  ipify-api:
    image: badsmoke/ipify-api
    build: .
    ports:
      - "3000:3000"
    environment:
      LOG_IPS: "true"
```

Start it with `docker compose up --build`. Stop a Compose deployment with:

```bash
docker compose down
```


## Questions?

Got a question? Please create a Github issue!
