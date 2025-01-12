# Gateor

Gateor is a lightweight, high-performance micro gateway built with Go.

## Features

- Authentication
- Rate limiting
- Reverse proxy
- Plugin-based architecture

## Installation

To install Gateor, you need to have Docker installed on your machine. Then, you can use the following commands to build and run the Docker container:

```bash
git clone https://github.com/yourusername/gateor.git
cd gateor
docker build -t gateor .
docker run -p 8080:8080 gateor
```


## Configuration

Gateor can be configured using YAML files in the `services` directory. Here is an example configuration for a service:

```yaml
name: demo-proxy
basepath: /demo
stripBasepath: true
target: 
  host: https://dummyapi.online
rateLimit: 5
```

## Plugin Architecture

Gateor uses a plugin-based architecture to handle request processing. Plugins can be chained together to create complex processing pipelines. The following plugins are available:

- `LeakyBucketRateLimit`: Limits the rate of incoming requests.
- `JwtAuthenticator`: Authenticates requests using JWT tokens.

## Contributing

We welcome contributions to Gateor! Please fork the repository and submit pull requests.
