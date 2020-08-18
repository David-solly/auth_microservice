# auth_microservice
Handles jwt provisioning, and verification. 
REQUIRED -
  - Redis instance to be running.
  
Client exposes a REST API with JSON payloads.
Server is a gRPC server that registers with hashicorp's consul for health checks

implementation of the [gokit-microservice-starter](https://github.com/David-solly/gokit-microservice-starter)

## jwt-service
Handles basic jwt authentication. Issues an expiring token pair, silent refresh and delete functions. corresponding to a persistent login session and logout when done.
