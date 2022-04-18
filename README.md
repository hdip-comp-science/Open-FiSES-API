# Go-REST-API

General overview and what you need to run this Go project.

## About the Application

This project implements a REST API server with the Gorilla Mux router. Objective is to build a service that carries out CRUS operations on both bookings and document services. The base URL for exposing HTTP endpoints is http://localhost:4000. A Postgres databse runs locally with Docker. Making use of the `gorm` package to talk to the database. 


## Prerequisites
- Familiarity with Go and PostgreSQL.
- You have Go and PostgreSQL installed on your machine.

## Usage

__Clone__ the repository to your machine: <br>
 ``` 
 git clone https://github.com/hdip-comp-science/Open-FiSES-API.git 
 cd Your_Repo_Directory 
 ```

If you have issue with dependecies, trying the following: <br>
`go mod init github.com/Open-FiSE/go-rest-api`

In order to run this API, you must have somewhere to store the document data. PostgreSQL is being used for storage. I chose to run Postgres locally with Docker. Use the following command:
`docker run --name postgres-db -e POSTGRES_PASSWORD=postgres -p 5433:5432 -d postgres` <br>
Check your running container using `docker ps` command. You should be running Postgres database locally now. 

Next, use the `Run` method to start the application: `go run server/main.go`. It will run through the connection logic and your API server should be accessible. Open Postman API platform or similar software to test the service endpoints. Example requests to the API service include: <br>


- __Creating__ a new document to a valid POST request `/document`
- __Updating__ a document in response to a valid PUT request `/document/{id}`
- __Deleting__ an existing document to a valid DELETE request `/document/{id}`
- __Getting__ an existing document based on ID `/document/{id}`, and fetching a __list__ of all documents `/documents`


## Project Package Imports

This is not an entire list of package imports utilised but one's I think are important.

- **Logrus** <br>
  `Logrus` is a structured logger for Go (golang), completely API compatible with the standard library logger. Logrus has the following characteristics:
  

    - Fully compatible with the log module of the golang standard library: logrus has six log levels: debug, info, warn, error, fatal and panic, which is a superset of the API of the log module of the golang standard library. If your project uses the standard library log module, it can be migrated to logrus at the lowest cost.
    - Extensible hook mechanism: allows users to distribute logs to any place through hook, such as local file system, standard output, logstash, elastic search or MQ, or define log content and format through hook.
    - Optional log output formats: logrus has two built-in log formats, jsonformatter and textformatter. If these two formats do not meet the requirements, you can implement the interface formatter by yourself to define your own log format.
    - Field mechanism: logrus encourages detailed and structured logging through the field mechanism, rather than logging through lengthy messages.
    - Logrus is a pluggable and structured log framework.
    <br>

- **Gorrila Mux** <br>
  The gorrilla/mux package [2] implements a request router and dispatcher for matching incoming requests to their respective handler. The name mux stands for “HTTP request multiplexer”. It is also compliant to Go’s default request handler signature `func (w http.ResponseWriter, r *http.Request)`, so the package can be mixed and machted with other HTTP libraries like middleware or exisiting applications. Use the go get command to install the package from GitHub like so: 
  &emsp;&emsp;&emsp;&emsp;&emsp;&emsp; ```go get -u github.com/gorilla/mux``` <br>

- **GORM** <br>
  GORM is a full featured Object Relationship Manager (ORM) [3]. It almost acts a broker between developers and the underlying database technology. They allow us to essentially work with object’s, much as we normally would and then save these objects without having to craft complex SQL statements [4].
  It supports database associations, preloading associated models, database transactions and much more. For more info check the [`documentation`](https://gorm.io/docs/) <br>

  <em>Side Note</em>: go-gorm has a prometheus package. Possible to monitor DB status with [`go-gorm/prometheus`](https://github.com/go-gorm/prometheus). <br>

- **encoding/json** <br>
  The `encoding/json` package is a standard library provided by GO.It is the most popular data format for sending and receiving data across the web. This package can automatically encode Go objects to JSON format. This is known as marshalling. The marshalling function accepts any Go data type and returns two values: a byte slice of the encoded JSON and an error. When dealing with user-genreated data and formatting it to JSON, it provides a highlevel of interoperability. Data stored in JSON makes it easy to exchange data and for database migration.<br>

- **net/http** <br>
  Package http provides HTTP client and server implementations.
  Get, Head, Post, and PostForm make HTTP (or HTTPS) requests. <br>

- **crypto/sha256** <br>
  The `crypto` package implements SHA224 and SHA256 hash algorithms as defined in FIPS 180-4.  This package plays an important role in the business logic of file uploading handler. The crypto package computes the hash value of the file and stores the value in the database. The hash value representation of the file is checked everytime athe upload endpoint is called.

- **os** <br>
  Package os provides a platform-independent interface to operating system functionality. The os interface is intended to be uniform across all operating systems. Features not generally available appear in the system-specific package syscall. 

## References
[1] [Golang File Upload](https://gabrieltanner.org/blog/golang-file-uploading)
[2] [Hash Checksum](https://yourbasic.org/golang/hash-md5-sha256-string-file/)
[3] [Go REST API Tutorial](https://tutorialedge.net/golang/creating-restful-api-with-golang/)
[4] [logrus documentation](https://pkg.go.dev/github.com/sirupsen/logrus#section-documentation)
[5] [gorrilla/mux](https://github.com/gorilla/mux)
[6] [GORM](https://github.com/go-gorm/gorm)
[7] [JSON](https://pkg.go.dev/encoding/json)
[8] [SHA256](https://pkg.go.dev/crypto/sha256#example-New)
[9] [Golang ORM Tutorial](https://tutorialedge.net/golang/golang-orm-tutorial/)
[10] [TutorialEdge.net](https://tutorialedge.net/)
