# RRS
## Preface
This is a project inspire by [Let's go](https://lets-go.alexedwards.net/) and [Let's go furthur](https://lets-go-further.alexedwards.net/) There author is [Alex Edwards](https://www.alexedwards.net/blog) He is a very brilliant man, I learn much from these two books, but those books require payment, So I can not just upload the codes to github, If U R new to Golang, I highly recommend you to buy those two books
## Introduction
When we learn a new language or a new tech, we need to read a lot of books, So this project is a way to manage the records of your reading, recommend, and interact with other readers, RRS is the backend of this system, I will write a simple CLI to interact with those APIs too([RRS_Client](https://github.com/rickj1ang/RRS_Client))
## Technique
A good way to learn is know what you will learn before leaning, So I will list the thchnology stack(also many skills maybe useful) of this project for your reference(this alse work as a catalog)
All in all to follow this project, every thing is free, RabbitMQ, MongoDB, Redis, and most important, I do not mean you deploy it by yourself, in local, but use some free tier cloud service, go to MongoDB, Redis official website, you can find the free tier, for RabbitMQ, I use [cloudamqp](https://www.cloudamqp.com/).
### Language
**Golang**: Go is becoming more and more popular nowadays, not only because it's clearity but powerful. It's official net/http package is a powerful weapon for backend development, and Goroutine is a very easy, and lightweight way to use multithreading for every coder.

**Makefile**: I also write a very simple Makefile for simplify my common instructions when develop, but the Makefile in [RRS_Client](https://github.com/rickj1ang/RRS_Client) is more standard.
### Web Framework
I use no third-party framwork for routing, middleware building, but [net/http](https://pkg.go.dev/net/http). I think we do not need to import more complexity for a simple project. But [httprouter](https://github.com/julienschmidt/httprouter) is a good package for first step, If you want to use some package. All in all, If you can use net/http well, you can use [Gin](https://github.com/gin-gonic/gin), [Gorilla](https://github.com/gorilla/mux) well too, because they are much easy but powerful
### RESTful
Same as [Let's go](https://lets-go.alexedwards.net/) and [Let's go furthur](https://lets-go-further.alexedwards.net/) this project follow RESRful to realize microservices architecture, rpc with [gRPC](https://grpc.io/) also a very popular microservices architecture in Golang but again we do not need to import more complexity for a simple project. RESTful is good enough for even for many big and commercial software
### MongoDB
I use MongoDB as my DataBase It's a very popular NoSQL DB, and U can follow me as well because I just use MongoDB atlas free tier for this project, which is enough. Go to. the official website login and then you can get your free account, I write my secret URI in a hard code way in my codespace, which is not good, but for simple
The reason of I choose MongoDB is the data I will put in is good structed. basiclly, it just custom go structs
### Redis
I use Redis as my Token cache, use MongoDB can do this well too, but Redis can set a expire time well, which is good for a token, token need to be expire in the future, as well, Redis is very fast, It can decrease the time user need to wait for the database reaction, So put something always use in Redis is a good choice
### Token
I simple use bearer token for my app to know who is the user, I store the hash of token as key, and user_id as value, in Redis as I said, when a user make a request with token header, app will go to Redis to get the user_id, and then, get the user_info from MongoDB use _id, which is very fase as index
### Log
I also custom some log, but not log to file, just out in the console, Which give us a unique style of output of our app 
### Error handle
In Go you must handle your error very carefully, sometimes, verbose, It is good choice for you to handle your error layer by layer, you do not handle error, when interact with database, but return the error to the handlers, the function call by endpoints, and you can collect your errors in a map so that you can output it clearly in json  
### Message Queue
I use RabbitMQ as my MQ, you can use any you like, this is a way I achive delay task, every time a user change a record and change the process of reading a certain book, app will put a message to a queue, and set a TTL, when it expire, the message will put into work queue by RabbitMQ exchange, I will have a Goroutine to listen this queue, and hanle the notify system, which will email the user to notify he/she it is the time to back to books  
### HealthCheck
I also use some package to read the cpu usage or ram useage, when you call healthcheck endpoint every time it will read this info, but it a lit bit expensive, do not call this endpoint too much