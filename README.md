# RRS
## Preface
This is a project inspire by [Let's go](https://lets-go.alexedwards.net/) and [Let's go furthur](https://lets-go-further.alexedwards.net/) There author is [Alex Edwards](https://www.alexedwards.net/blog) He is a very brilliant man, I learn much from these two books, but those books require payment, So I can not just upload the codes to github, If U R new to Golang, I highly recommend you to buy those two books
## Introduction
When we learn a new language or a new tech, we need to read a lot of books, So this project is a way to manage the records of your reading, recommend, and interact with other readers, RRS is the backend of this system, I will write a simple CLI to interact with those APIs too([RRS_Client](https://github.com/rickj1ang/RRS_Client))
## Technique
A good way to learn is know what you will learn before leaning, So I will list the thchnology stack(also many skills maybe useful) of this project for your reference(this alse work as a catalog)
### Language
**Golang**: Go is becoming more and more popular nowadays, not only because it's clearity but powerful. It's official net/http package is a powerful weapon for backend development, and Goroutine is a very easy, and lightweight way to use multithreading for every coder.

**Makefile**: I also write a very simple Makefile for simplify my common instructions when develop, but the Makefile in [RRS_Client](https://github.com/rickj1ang/RRS_Client) is more standard.
### Web Framework
I use no third-party framwork for routing, middleware building, but [net/http](https://pkg.go.dev/net/http). I think we do not need to import more complexity for a simple project. But [httprouter](https://github.com/julienschmidt/httprouter) is a good package for first step, If you want to use some package. All in all, If you can use net/http well, you can use [Gin](https://github.com/gin-gonic/gin), [Gorilla](https://github.com/gorilla/mux) well too, because they are much easy but powerful
### RESTful
Same as [Let's go](https://lets-go.alexedwards.net/) and [Let's go furthur](https://lets-go-further.alexedwards.net/) this project follow RESRful to realize microservices architecture, rpc with [gRPC](https://grpc.io/) also a very popular microservices architecture in Golang but again we do not need to import more complexity for a simple project. RESTful is good enough for even for many big and commercial software
### MongoDB
I use MongoDB as my DataBase 
### Redis
### Token
### Log
### Error handle
### Message Queue
### Authentication
### HealthCheck