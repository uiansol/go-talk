# go-talk
A chat server writen in Golang

This project is a demonstration of Golang usage on 
back-end development. It uses some important concepts while
focus on simplicity to be used as reference to take code pieces.

**What you will find here**

* Routing with [chi](https://github.com/go-chi/chi)

* Data base usage with [GORM](https://gorm.io/)

* Log server messages to file

* Use [JWT](https://en.wikipedia.org/wiki/JSON_Web_Token) Token when login

* Use [crypto](https://golang.org/pkg/crypto/) for password storage

* Deploy with [Docker](https://www.docker.com/)

* Use [RabbitMQ](https://www.rabbitmq.com/) to publish messages

* Unit tests parallelized.

### How to run

Build and run with docker-compose
`docker-compose up`

Visit on a browser
`localhost:8080/`

Run tests (-v to see the parallelism)
`go test -v`

### Future improvements and engineering discussions

* Sqlite3 is a good data base for small or demo purposes. As we are using GORM it is enough for now. For a better architecture we use one container with PostgreSQL or MySQL. In code the modification will be some env variables to store the database credentials and dsn (data source name).

* We could move to a web sockets model instead of using polling. The primary difference here is that instead of users constantly requesting messages, we would push new messages to the users via web sockets. For now it was intended to reduce work time on front end to demonstrate more features on back end with Go.

* Tests for the api and end to end is a plus. But it would require data mocking and a test database that is outside the scope of this project. The few we have are enough to demonstrate how to deal with unit tests in Go (using tables and parallelizing it).

* The consumer for messaging are in another repository, just to illustrate a microservice feeling.