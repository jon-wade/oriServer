# oriServer

## Installation

## Overview

The `oriServer` is a gRPC server that responds to `oriClient`, a component part of the Ori tech test. The application is
containerised using Docker, and contains a CI/CD pipeline using GitHub Actions that runs the unit test suite, building a
new Docker image and pushing it up to the Docker Hub, ready for deployment into a Kubernetes cluster, after every push
to the remote repo.

Kubernetes manifests have been included, specifically a `Deployment` that creates a single pod and a `Service` that
exposes the application on its default port `50051` via a cloud-based load balancer. The manifests have been tested on
an AWS EKS cluster and are reachable using the `oriClient` correctly configured. 

As there is only one external configuration variable, the `ORI_PORT` number, the env var is set in
the cluster via the `Deployment`.

## Pre-requisites

The `oriServer` runs within a Docker container and makes use of Docker-compose, therefore both need to be installed
before attempting to spin up the application. 

Docker installation instructions are available here: 
https://docs.docker.com/install/

Docker-compose installation instructions here: https://docs.docker.com/compose/install/

`go` must be installed, version 1.13.x, to run unit tests. For installation details, please see here:
https://golang.org/doc/install

## Installation
 1. Clone the repo with `git clone https://github.com/jon-wade/oriServer.git`
 2. `cd oriServer` to change directories to the newly cloned repo
 3. `go mod download` to install the dependencies
 
It is now possible to run the unit tests directly as per the testing section below.

### Spin-Up
#### Development
`ORI_PORT=50051 docker-compose -f docker-compose.dev.yml up` to spin up the application on `localhost:50051`, with a 
watcher application to build and recompile on file changes. This will either pull or build the base image on first use.
#### Production
`ORI_PORT=50051 docker-compose up` to spin up the application on `localhost:50051` without file watching etc.

### Tear-Down
`ORI_PORT=50051 docker-compose down` to tear down either environment.

### Running tests
Unit tests can be run using `go test ./...`

## Documentation
The following details the responses to the questions raised in the documentation of the Ori tech test.

### Alignment to 12factor app best practices

Reference: https://12factor.net/

* **Codebase**: tracked in Git revision control, small incremental commits and automated deployments
* **Dependencies**: explicitly tracked and isolated using go modules.
* **Config**: configuration stored in environment variables.
* **Backing services**: the server currently has no backing services.
* **Build, release, run**: currently the GitHub Actions pipeline runs unit tests and then as a follow-up step, builds a
Docker image and pushes it to the DockerHub. Follow-up steps would then update environment variables and the Docker
image running in the K8 cluster.
* **Processes**: the app is stateless and containerised, and can be executed as one or more processes in its environment.
* **Port Binding**: the app is exposed externally via a K8 service, which maps an external port to the internal cluster port.
* **Concurrency**: the app can be scaled by spinning up more containers, as processes above.
* **Disposability**: in the development environment a simple one line command is all that is required to spin up or tear down
the app. In a production K8, pods running containers are mortal and can be scaled up or down by a single CLI command.
* **Dev/Prod parity**: The development environment runs within the Docker container that will run on production isolating
local development environment configuration issues from the application runtime.
* **Logs**: In production, logs would likely be exported from the cluster as an event stream using something like 
FluentD and collected in a tool to query the streams, for instance CloudWatch.
* **Admin processes**: Docker and Docker-compose files are committed to the repo and can be run as one-off processes.

### Best cloud native understanding

Reference: https://thenewstack.io/10-key-attributes-of-cloud-native-applications/

1. **Packaged as lightweight containers**: Docker has been used to create the container image, utilising a multi-stage
build to make container as lightweight as possible.
2. **Developed with best-of-breed languages and frameworks**: The application has been built in Go, a highly-performant
statically typed language and utilises gRPC, the de facto standard for microservice communication.
3. **Designed as loosely coupled microservices**: the two microservices that make up the Ori system, the `oriServer` and
`oriClient` are loosely coupled, existing as unique applications sharing only one dependency, the `ori.pg.go` protobuf
configuration file, which is imported as a go module dependency into `oriClient`.
4. **Centered around APIs for interaction and collaboration**: both applications are centred around the gRPC protobuf
definition to simplify interaction and collaboration.
5. **Architected with a clean separation of stateless and stateful services**: both applications are currently stateless.
6. **Isolated from server and operating system dependencies**: the `oriServer` is isolated via its Docker container, the
`oriClient` via the cross-platform compilation capabilities of the `Go` language.
7. **Deployed on self-service, elastic, cloud infrastructure**: the `oriServer` application is intended to be released
into a K8 cluster, which provides automatic scaling and fault-tolerance. It has been tested on an AWS EKS cluster, a
cloud-based implementation of K8.
8. **Managed through agile DevOps processes**: the `oriServer` has a CI/CD pipeline that runs unit tests and builds
and pushes Docker images automatically, on each push to the remote repo.
9. **Automated capabilities**: as mentioned, CI/CD and automation are built in as core to the system.
10. **Defined, policy-driven resource allocation**: policy implementation is made easy via automated deployments and the
pipeline and K8 configuration files.

### Expansion to incorporate an eventstore

To add an eventstore to the application, initially I would consider which _events_ would need to be stored. As the
application is very simple, with only two methods, the obvious events are client requests and server responses. Each
would require a data structure to store details, so for instance for client requests, we might create a struct such as
`calculationRequested` with fields such as `in` representing the inputs, `ip` the address of the requesting 
service, `ts` a timestamp of the request and `method`, the type of calculation requested.

These events need to be stored somewhere, so some form of data store would need to be added to the application. A table
in a cloud-hosted backed-up relational database would suffice for this purpose, which would store the struct fields, 
along with a `created_at` timestamp as well as details of the application's container, to be able to identify which 
instance of the application created the record. The application would require code to be added to allow access to that 
resource, which would also require some amendments to the external configuration to hold the connection credentials.

Writes to this store should happen concurrently with the main processing logic using `go routines` to ensure that 
performance is not impacted. Its also likely best to cache events locally in memory for a period and then write in 
batches, to reduce network latency impacts. This approach does introduce the risk of data loss in the case of a failed 
application, so this risk would need to be balanced against overall performance considerations.

Next, I would turn my attention to other components of an eventstore. Firstly, do we require an `aggregate` to hold the
current state of the app? In the case of our app, we may want to hold a running total of the number and types of
calculation performed by the server at any point in time. Again, in this case a new data structure could be created to
store that data such as `runningTotal` with integer fields such as `summation` and `factorial` to hold the totals. We
would then need to write `calculators` that implement the logic to update the totals, which would likely take the form of
methods added to the event structs. Finally, does the application require any `reactors`, processes that react to certain
events? For example, do we require notifications when certain client ips make requests? If so, adding additional logic
to generate those notification could be added, perhaps to the event dispatcher.

## Access from outside the cluster
The application could be made available for access outside the cluster via a Kubernetes `Service` which simply maps the 
internal cluster port to an external port. On AWS, for instance, its then possible to set the the service `type` to 
`LoadBalancer` and an external load balancer will then be automatically configured which will provide a resolvable 
address to call the service over the internet. There are other methods of configuring access, including `NodePort`, 
`ExternalName` types or making use of an `Ingress` manifest. An example service manifest, `k8/oriserver-service.yaml` is
available in this repo.

