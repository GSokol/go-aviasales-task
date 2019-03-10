Aviasales go task
=================

This is aviasales go task

Usage
-----

### Unit tests

Run unit tests locally (project should be located at `$GOPATH/src/GSokol/go-aviasales-task`):

```bash
$ make test
```

### K8S

Deploy to k8s: (you should have k8s custer runninig with ingress controller installed)

```bash
$ make -f devenv/Makefile deploy-k8s
```

NOTE: to ship images you should setup your registry with `DOCKER_REGISTRY` environment variable and set `SHOULD_PUSH=1` (see devenv/Makefile).
NOTE: if you use minikube just do `eval $(minikube docker-env)`

### Docker-compose

To deploy locally with compose of 2 containers (redis + service):

```bash
$ docker-compose up -d
```

### Pure Docker file

To run using pure docker, first build the image:

```bash
$ docker build -t go-aviasales-task .
```

Then run image:

```bash
$ docker run -d \
  -p 8080:8080 \
  -v ./etc:/etc/go-aviasales-task \
  -e AV_CONFIG_PATH=/etc/go-aviasales-task/docker.json \
  go-aviasales-task
```
