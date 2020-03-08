#oriServer

## Operation

### Spin-Up
####Development
`ORI_PORT=50051 docker-compose -f docker-compose.dev.yml up` to spin up the application on port 50051, with a watcher 
application to build and recompile on file changes.

####Production
`ORI_PORT=50051 docker-compose up` to spin up the application on port 50051.

### Tear-Down
`ORI_PORT=50051 docker-compose down` to tear down either environment.

## Documentation
The following details the responses to the questions raised in the documentation of the Ori tech test.

### Alignment to 12factor app best practices

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

TODO

### Expansion to incorporate an eventstore

To add an eventstore to the application, initially I would consider which _events_ would need to be stored. As the
application is very simple, with only two methods, the obvious events are client requests and server responses. Each
would require a data structure to store details, so for instance for client requests, we might create a struct such as
`calculationRequested` with fields such as `in` representing the inputs, `ip` the address of the requesting 
service, `ts` a timestamp of the request and `method`, the type of calculation requested.

These events need to be stored somewhere, so some form of data store would need to be added to the application. A table
in a cloud-hosted backed-up relational database would suffice for this purpose, which would store the struct fields, 
along with a `created_at` timestamp as well as details of the applications container, to be able to identify which 
instance of the application created the record. The application would require code to be added to allow access to that 
resource, which would also require some amendments to the external configuration to hold the connection credentials.

Writes to this store should happen concurrently with the main processing logic using `go routines` to ensure that 
performance is not impacted. Decisions would also need to be made on the frequency of writes to this database. For high 
performance, its likely best to cache events locally in memory for a period and then write in batches, to reduce network
delays. This approach does introduce the risk of data loss in the case of a failed application, so this risk would
need to be balanced against overall performance considerations.

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

