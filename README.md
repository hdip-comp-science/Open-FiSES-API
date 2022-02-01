# Create-GO-API

General overview and what you need to run this Go project.

## Current Status
Objective is to build a service that handles documents, get, add, update and delete. 
This project implements a REST HTTP Endpoints with the Gorilla Mux router. A Postgres databse runs locally with Docker. Making use of the `gorm` package to talk to the database. 


## Usage
In order to run this API, you must have somewhere to store the document data. PostgreSQL is being used for storage. I chose to run Postgres locally with Docker. Use the following command:
`docker run --name postgres-db -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres` <br>
Check your running container using `docker ps` command. You should be running Postgres db locally now. For now, export the database environment variables from the CLI:
```
 export DB_USERNAME=postgres
 export DB_PASSWORD=postgres
 export DB_HOST=localhost
 export DB_TABLE=postgres
 export DB_PORT=5432
 export DB_DB=postgres
```
Next, run the application: `go run server/main.go`. It will run through the connection logic and your API server should be accessible. Open Postman to test the document service endpoints.


## Project Package Imports

This is not an entire list of package imports but one's I thought are worth mentioning.

- **GLOG** <br>
  The `glog` package implements a pure Go implementation of levelled logs. It provides the functions Info, Warning, Error, Fatal, plus formatting variants such as Infof, Warningf etc. Checkout glog repository[1] for more details.<br>

- **Gorrila Mux** <br>
  The gorrilla/mux package [2] implements a request router and dispatcher for matching incoming requests to their respective handler. The name mux stands for “HTTP request multiplexer”. It is also compliant to Go’s default request handler signature `func (w http.ResponseWriter, r *http.Request)`, so the package can be mixed and machted with other HTTP libraries like middleware or exisiting applications. Use the go get command to install the package from GitHub like so: 
  &emsp;&emsp;&emsp;&emsp;&emsp;&emsp; ```go get -u github.com/gorilla/mux``` <br>

- **GORM** <br>
  GORM is a full featured Object Relationship Manager (ORM) [3]. It almost acts a broker between developers and the underlying database technology. They allow us to essentially work with object’s, much as we normally would and then save these objects without having to craft complex SQL statements [4].
  It supports database associations, preloading associated models, database transactions and much more. For more info check the [`documentation`](https://gorm.io/docs/) <br>

  <em>Side Note</em>: go-gorm has a prometheus package. Possible to monitor DB status with [`go-gorm/prometheus`](https://github.com/go-gorm/prometheus). <br>

- **net/http** <br>
  Package http provides HTTP client and server implementations.
  Get, Head, Post, and PostForm make HTTP (or HTTPS) requests: <br>

- **os** <br>
  Package os provides a platform-independent interface to operating system functionality. The os interface is intended to be uniform across all operating systems. Features not generally available appear in the system-specific package syscall. 

## References

[1] [glog repository](https://github.com/golang/glog)
[2] [gorrilla/mux](https://github.com/gorilla/mux)
[3] [GORM`](https://github.com/go-gorm/gorm)
[4] [Golang ORM Tutorial](https://tutorialedge.net/golang/golang-orm-tutorial/)
[5] [TutorialEdge.net](https://tutorialedge.net/)
