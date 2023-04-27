# Book Tracker Service

A golang based microservice that manages the reading activity of users that provides the below functionalities :
- Add a book to the reading list
- Update the book(Example: Set the status to IN PROGRESS, Bookmark a page)
- List books(sorted by status or title)
- Fetch a specific book
- Delete the book(it is a soft delete - meaning the Front End would call the Update endpoint with active="false")

## Structure
The structure of the project is following the architecture proposed by Robert C. Martin - [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

The layers proposed by this architecture are inside the internal folder:

* adapter (Interface Adapter)
* entity (Entities - Enterprise Business Roles)
* framework (Framework and Drivers)
* service (Use Cases - Application Business Rules)

```
|-- api
|   |-- openapi.json
|-- build
|-- cmd
|   |-- microservice
|-- internal
|   |-- adapter
        |-- webserver
            |-- probes
            |-- swagger
|   |-- entity
|   |-- framework
        |-- database
|   |-- service
|-- kube
|-- tests
|   |-- results
|-- .gitignore
|-- CHANGELOG.md
|-- docker-compose.yml
|-- Dockerfile
|-- go.mod
|-- go.sum
|-- Makefile
|-- README.md

```

## Prerequisites
```
Go version >= 1.20
```

## Install
    $ git clone https://github.com/AnushaSankaranarayanan/book-tracker-service.git
    $ cd book-tracker-service

## Build, Tests and Coverage

#### Clean, test and build

- \$ make  all

#### Build

- \$ make build

#### Run tests

- \$ make  test

#### Run tests with coverage

- \$ make cover

---

## Documentation
The module contains the documentation from code, so all comments can be generated in HTML.

**Get the godoc**
```
go get golang.org/x/tools/cmd/godoc
```

**Run the command below**
```
godoc -http=:6060
```

**View docs**

Open godocs [here](http://localhost:6060/pkg/github.com/anushasankaranarayanan/book-tracker-service/)

## Running the service locally

*Note : Please have a running couchbase server before proceeding to the below section. Refer to section Couchbase Prerequisites for the initial setup.

### Non containerized
Create .env file and set the values for below properties. (Sample values given below)
```
LOG_LEVEL=INFO
SERVER_PORT=9000
NAME=book-tracker-service
COUCHBASE_HOST=localhost:8091
COUCHBASE_BUCKET=reading-list
COUCHBASE_USER=<username>
COUCHBASE_PASSWORD=<password>
ENABLE_DB_VERBOSE_LOGGING=false

```
Navigate to directory:
```
cd cmd/microservice

```
Run `go run -tags real main.go` . Issue requests to : `http://localhost:9000`. Refer to `http://localhost:9000/api/v1/openapi` for the list of endpoints supported.

### Docker compose
#### Install [Docker](https://www.docker.com/) on your system.

* [Install instructions](https://docs.docker.com/installation/mac/) for Mac OS X
* [Install instructions](https://docs.docker.com/installation/ubuntulinux/) for Ubuntu Linux
* [Install instructions](https://docs.docker.com/installation/) for other platforms

#### Install [Docker Compose](http://docs.docker.com/compose/) on your system.

* Python/pip: `sudo pip install -U docker-compose`
* Other: ``curl -L https://github.com/docker/compose/releases/download/1.1.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose; chmod +x /usr/local/bin/docker-compose``

#### Run `make build` command on the project root directory.

#### Edit `docker-compose.yml` with appropriate environment variables.

#### Run the below command to bring up the service.
```
 docker-compose up --build
```
#### Once the application is up - issue requests to : `http://localhost:9000`.

#### To stop the containers
 ```
docker-compose stop
```
### Kubernetes cluster

#### Ensure that a kubernetes cluster is up and running.For minikube, please refer [here](https://minikube.sigs.k8s.io/docs/start/)

####  Ensure that `kubectl` has access to the running cluster.

####  Run the below command
```
cd kube
./start.sh
```

The service should be up and running.Get the URL of the service by running:
```
minikube service book-tracker-service --url
```
and issue requests
#### To stop and remove the deployments from kubernetes , run
```
/kube/clean-up.sh
```

## Sample responses
```
# Liveness probe
curl --location --request GET 'http://localhost:9000/api/v1/probes/liveness'

{ "name": "book-tracker-service" }

# Add a book - error scenario

curl --location 'http://localhost:9000/api/v1/book' \
--data '{
    "isbn": "978-1-60309-527-3",
    "title": "But You Have Friends",
    "author": "Emilia McKenzie",
    "genre": "Adventure"
}'

{
    "code": 400,
    "status": "Bad Request",
    "message": "Key: 'Book.ISBN' Error:Field validation for 'ISBN' failed on the 'required' tag"
}

# Add a book - success scenario

curl --location 'http://localhost:9000/api/v1/book' \
--data '{
    "isbn": "978-1-60309-527-3",
    "title": "But You Have Friends",
    "author": "Emilia McKenzie",
    "genre": "Adventure"
}'

{
    "code": 200,
    "status": "OK",
    "message": "book creation successful"
}

# List books

curl --location 'http://localhost:9000/api/v1/book/'

{
    "code": 200,
    "status": "OK",
    "message": "books retrieval successful",
    "count": 8,
    "books": [
        {
            "isbn": "978-1-60309-038-4",
            "title": "Essex County",
            "author": "Jeff Lemire",
            "genre": "Thriller",
            "updated": 1682514188,
            "started": 1682513807
        },
        {
            "isbn": "978-1-60309-084-1",
            "title": "Does Something",
            "author": "James Kochalka",
            "genre": "Thriller"
        },
        {
            "isbn": "978-1-60309-329-3",
            "title": "The Tempest",
            "author": "Alan Moore",
            "genre": "Mystery"
        },
        {
            "isbn": "978-1-60309-469-6",
            "title": "From Hell",
            "author": "Eddie Campbell",
            "genre": "Horror"
        },
        {
            "isbn": "978-1-60309-481-8",
            "title": "Parenthesis",
            "author": "Lodie Durand",
            "genre": "Horror"
        },
        {
            "isbn": "978-1-60309-504-4",
            "title": "Glork Patrol Takes a Bath",
            "author": "James Kochalka",
            "genre": "Mystery"
        },
        {
            "isbn": "978-1-60309-513-6",
            "title": "Doughnuts and Doom",
            "author": "Balazs Lorinczi",
            "genre": "Mystery"
        },
        {
            "isbn": "978-1-60309-527-3",
            "title": "But You Have Friends",
            "author": "Emilia McKenzie",
            "genre": "Adventure",
            "created": 1682514622,
            "updated": 1682514622,
            "created_by": "SYSTEM",
            "updated_by": "SYSTEM",
        }
    ]
}
# List books - sorted by title
curl --location 'http://localhost:9000/api/v1/book?sort=title'

{
    "code": 200,
    "status": "OK",
    "message": "books retrieval successful",
    "count": 2,
    "books": [
        {
            "isbn": "978-1-60309-527-3",
            "title": "But You Have Friends",
            "author": "Emilia McKenzie",
            "genre": "Adventure",
            "status": "FINISHED",
            "created": 1682585792,
            "updated": 1682588831,
            "created_by": "SYSTEM",
            "updated_by": "SYSTEM",
            "active": "true"
        },
        {
            "isbn": "978-1-60309-469-6",
            "title": "From Hell",
            "author": "Eddie Campbell",
            "genre": "Horror",
            "status": "IN PROGRESS",
            "updated": 1682588775,
            "active": "true"
        }
    ]
}

# Get a book - error scenario(invalid id)

curl --location 'http://localhost:9000/api/v1/book/bla'
{
    "code": 404,
    "status": "Not Found",
    "message": "book with id bla not found"
}

# Get a book - success scenario

curl --location 'http://localhost:9000/api/v1/book/978-1-60309-038-4'
{
    "code": 200,
    "status": "OK",
    "message": "book retrieval successful",
    "book": {
        "isbn": "978-1-60309-038-4",
        "title": "Essex County",
        "author": "Jeff Lemire",
        "genre": "Thriller",
        "updated": 1682514188,
        "started": 1682513807,
        "active": "false"
    }
}

# Update a book - error scenario(invalid id)

curl --location --request PUT 'http://localhost:9000/api/v1/book' \
--data '{
    "isbn": "bla",
    "title": "Essex County",
    "author": "Jeff Lemire",
    "genre": "Thriller",
    "active": "false",
      "started": 1682513807

}'

{
    "code": 404,
    "status": "Not Found",
    "message": "book with id bla not found"
}

# Update a book - success scenario

curl --location --request PUT 'http://localhost:9000/api/v1/book' \
--data '{
    "isbn": "978-1-60309-038-4",
    "title": "Essex County",
    "author": "Jeff Lemire",
    "genre": "Thriller",
    "started": 1682513807,
    "bookmark": 100
}'
{
    "code": 200,
    "status": "OK",
    "message": "book updated successfully"
}

```
## Known caveats
* Swagger assets are included in the service. Moving that to a common module would be a sensible choice
* The service doesn't provide a separate endpoint for DELETE service. Calling the UDPATE endpoint with active="false" is advised for such cases
* Couchbase is used as DB here . This could be changed to any DB after an elaborate internal discussion with the team
* Build tags(fake and real) are used for unit testing. Alternatively , we could use mocks
* The microservice is not guarded currently. Ideally this could be achieved using OAuth or any other mechanisms per the team standards
* Logrus is being used for logging. This could be moved ad used as a middleware to prevent initializing in multiple places
* sonar.properties file could be included
* List endpoint is not paginated . It only returns a count . This could be modified to accept limit parameter to aid in pagination
* Books cannot be exported to the pantry basket using this service since there is no equivalent library for Go
* Environment Variable COUCHBASE_PASSWORD is set as plain text, this SHOULD be moved to a secret.
* Sort is on ascending order. This could be driven by a query parameter. The sorting is done at the service rather than database queries(as sort queries can be a bit expensive) 

## Feature Improvements 
* The data model has a field called "bookmark" which can be used to track the progress of the user. It can be set when calling the UPDATE endpoint. The user could be directly taken to the page when he/she selects the book from the UI.
* The Front end can use the timestamps(start/end) returned from the LIST endpoint to show a dashboard / graph to the user showing weekly reading times,
* The service supports multi tenancy by default by leveraging couchbase scopes and collections[Documentation here](https://docs.couchbase.com/server/current/learn/data/scopes-and-collections.html).So in the future reading lists for a family can be added without much code changes
## Couchbase Prerequisites
* Create a bucket called `reading-list` in the couchbase cluster . [Refer here for more details](https://docs.couchbase.com/server/current/manage/manage-buckets/create-bucket.html)
* Run the below queries in the couchbase server for the initial setup:
```
CREATE COLLECTION `reading-list`.`_default`.book
CREATE PRIMARY INDEX primary_index_book on `reading-list`.`_default`.book;

```
