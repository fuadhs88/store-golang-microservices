# 🧺 Store Golang Microservices

[![CI](https://github.com/mehdihadeli/store-golang-microservices/actions/workflows/ci.yml/badge.svg?branch=main&style=flat-square)](https://github.com/mehdihadeli/store-golang-microservices/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mehdihadeli/store-golang-microservice-sample)](https://goreportcard.com/report/github.com/mehdihadeli/store-golang-microservice-sample)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod&style=flat-square&color=ff69b4)](https://gitpod.io/#https://github.com/mehdihadeli/store-golang-microservices-sample)

![store-golang-microservices](https://socialify.git.ci/mehdihadeli/store-golang-microservices/image?description=1&descriptionEditable=%F0%9F%A7%BA%20Practical%20store%20microservices%2C%20built%20with%20Golang%2C%20DDD%2C%20CQRS%2C%20Event%20Sourcing%2C%20Vertical%20Slice%20Architecture%2C%20and%20the%20latest%20technologies.&font=Inter&forks=1&issues=1&language=1&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Light)

> A practical store sample, built with Golang and different software architecture and technologies like **Microservices Architecture**, **Vertical Slice Architecture** , **CQRS Pattern**, **Domain Driven Design (DDD)**, **Event Driven Architecture**. For communication between independent services, We use asynchronous messaging with using kafka, and sometimes we use synchronous communication for real-time communications with using REST and gRPC calls.

🌀 This Application is `in-progress` and I will add new features and thecnologies over time. 


## Technologies - Libraries

- ✔️ **[`labstack/echo`](https://github.com/labstack/echo)** - High performance, minimalist Go web framework
- ✔️ **[`segmentio/kafka-go`](https://github.com/segmentio/kafka-go)** - Kafka library in Go
- ✔️ **[`grpc/grpc-go`](https://github.com/grpc/grpc-go)** - The Go language implementation of gRPC. HTTP/2 based RPC
- ✔️ **[`jackc/pgx`](https://github.com/jackc/pgx)** - PostgreSQL driver and toolkit for Go
- ✔️ **[`swaggo/swag`](https://github.com/swaggo/swag)** - Automatically generate RESTful API documentation with Swagger 2.0 for Go.
- ✔️ **[`jaegertracing/jaeger-client-go`](https://github.com/jaegertracing/jaeger-client-go)** - Distributed tracing for go (will replace with open telemetry soon)
- ✔️ **[`prometheus/client_golang`](github.com/prometheus/client_golang)** - Prometheus instrumentation library for Go applications
- ✔️ **[`mongodb/mongo-go-driver`](https://github.com/mongodb/mongo-go-driver)** - The Go driver for MongoDB
- ✔️ **[`go-redis/redis`](https://github.com/go-redis/redis)** - Type-safe Redis client for Golang
- ✔️ **[`go-gorm/gorm`](https://github.com/go-gorm/gorm)** - The fantastic ORM library for Golang, aims to be developer friendly
- ✔️ **[`mschneider82/problem`](https://github.com/mschneider82/problem)** - A go library that implements application/problem+json and application/problem+xml
- ✔️ **[`go-playground/validator`](https://github.com/go-playground/validator)** - Go Struct and Field validation, including Cross Field, Cross Struct, Map, Slice and Array diving

## Plan
> This project is in progress, New features will be added over time.

High-level plan is represented in the table

| Feature | Status |
| ------- | ------ |
| Building Blocks | In Progress 👷‍ |
| Catalog Service | In Progress 👷‍ |
| API Gateway | Not Started 🚩 |
| Identity Service | Not Started 🚩 |
| Customer Service | Not Started 🚩 |
| Order Service |  In Progress 👷‍|

## Application Structure

In this project I used [vertical slice architecture](https://jimmybogard.com/vertical-slice-architecture/) or [Restructuring to a Vertical Slice Architecture](https://codeopinion.com/restructuring-to-a-vertical-slice-architecture/) also I used [feature folder structure](http://www.kamilgrzybek.com/design/feature-folders/) in this project.

- We treat each request as a distinct use case or slice, encapsulating and grouping all concerns from front-end to back.
- When We adding or changing a feature in an application in n-tire architecture, we are typically touching many different "layers" in an application. we are changing the user interface, adding fields to models, modifying validation, and so on. Instead of coupling across a layer, we couple vertically along a slice and each change affects only one slice.
- We `Minimize coupling` `between slices`, and `maximize coupling` `in a slice`.
- With this approach, each of our vertical slices can decide for itself how to best fulfill the request. New features only add code, we're not changing shared code and worrying about side effects. For implementing vertical slice architecture using cqrs pattern is a good match.

![](./assets/vertical-slice-architecture.jpg)


![](./assets/vsa2.png)

Also here I used [CQRS](https://www.eventecommerce.com/cqrs-pattern) for decompose my features to very small parts that makes our application:

- maximize performance, scalability and simplicity.
- adding new feature to this mechanism is very easy without any breaking change in other part of our codes. New features only add code, we're not changing shared code and worrying about side effects.
- easy to maintain and any changes only affect on one command or query (or a slice) and avoid any breaking changes on other parts
- it gives us better separation of concerns and cross cutting concern (with help of MediatR behavior pipelines) in our code instead of a big service class for doing a lot of things.

With using [CQRS](https://event-driven.io/en/cqrs_facts_and_myths_explained/), our code will be more aligned with [SOLID principles](https://en.wikipedia.org/wiki/SOLID), especially with:

- [Single Responsibility](https://en.wikipedia.org/wiki/Single-responsibility_principle) rule - because logic responsible for a given operation is enclosed in its own type.
- [Open-Closed](https://en.wikipedia.org/wiki/Open%E2%80%93closed_principle) rule - because to add new operation you don’t need to edit any of the existing types, instead you need to add a new file with a new type representing that operation.

Here instead of some [Technical Splitting](http://www.kamilgrzybek.com/design/feature-folders/) for example a folder or layer for our `services`, `controllers` and `data models` which increase dependencies between our technical splitting and also jump between layers or folders, We cut each business functionality into some vertical slices, and inner each of these slices we have [Technical Folders Structure](http://www.kamilgrzybek.com/design/feature-folders/) specific to that feature (command, handlers, infrastructure, repository, controllers, data models, ...).

Usually, when we work on a given functionality we need some technical things for example:

- API endpoint (Controller)
- Request Input (Dto)
- Request Output (Dto)
- Some class to handle Request, For example Command and Command Handler or Query and Query Handler
- Data Model

Now we could all of these things beside each other and it decrease jumping and dependencies between some layers or folders.

Keeping such a split works great with CQRS. It segregates our operations and slices the application code vertically instead of horizontally. In Our CQRS pattern each command/query handler is a separate slice. This is where you can reduce coupling between layers. Each handler can be a separated code unit, even copy/pasted. Thanks to that, we can tune down the specific method to not follow general conventions (e.g. use custom SQL query or even different storage). In a traditional layered architecture, when we change the core generic mechanism in one layer, it can impact all methods.

### Catalogs Service Structure

![](./assets/catalogs-service.png)

### High Level Structure

![](./assets/high-level.png)


## References

- [https://github.com/AleksK1NG/Go-CQRS-Kafka-gRPC-Microservices](https://github.com/AleksK1NG/Go-CQRS-Kafka-gRPC-Microservices)
- [https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
- [https://github.com/kgrzybek/modular-monolith-with-ddd](https://github.com/kgrzybek/modular-monolith-with-ddd)
- [https://github.com/jbogard/ContosoUniversityDotNetCore-Pages](https://github.com/jbogard/ContosoUniversityDotNetCore-Pages)
