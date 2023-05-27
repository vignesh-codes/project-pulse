
# Pulse Service

Go --> logfiles --> fluent-bit --> go-fluent-helper-aggregator --> redis --> go-fluent-helper-aggregator --> postgres

Ideally you might want to remove go-fluent-helper-aggregator and identify a good log aggregator with appropriate INPUT and OUTPUT plugins.

A Go http based log aggregator that collects logs via api from various internal serices and puts them in log files.  
Fluent-bit tails those logs and pushes to another go-api service puts them in redis streams.

A go-aggregator service periodically checks the redis streams. And if data is present, it captures at a batch of 1k items and batchwrites to 
postgres.

### To start the project
```
cd pulse-service
docker build -t pulse .

cd aggregator-service
docker build -t agg-v1 .

cd test
docker build -t test-client .

cd fluent-bit
docker built -t fluent .

docker compose up redis postgres pulse agg fluent -d

docker compose up test 
```

```
In case if you dont want docker compose
docker run --network host --privileged -v /Users/logs:/logs -dit fluent-v1

docker run --network host -dit agg-v1

docker run --network host -v /Users/logs:/go/logs -dit pulse

docker run --name redis -p 6379:6379 -d redis

docker run --name postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -p 4432:5432 -d postgres
```

```
.
└── project-pulse
    ├── README.md
    ├── aggregator-service
    │   ├── dockerfile
    │   ├── go.mod
    │   ├── go.sum
    │   ├── main.go
    │   ├── models.go
    │   └── pushclient.go
    ├── apps
    │   ├── controller
    │   │   ├── client
    │   │   │   └── v1
    │   │   └── private
    │   │       └── v1
    │   │           └── v1.internal.controller.go
    │   ├── dao
    │   │   ├── client
    │   │   │   └── v1
    │   │   └── private
    │   │       └── v1
    │   │           └── v1.internal.dao.activitylogger.go
    │   ├── repository
    │   │   ├── adapter
    │   │   │   ├── adapter.go
    │   │   │   ├── psql.go
    │   │   │   └── redis.go
    │   │   └── instance
    │   │       └── instance.go
    │   ├── routes
    │   │   ├── client
    │   │   │   └── v1.client.routes.go
    │   │   ├── config.go
    │   │   ├── internal
    │   │   │   └── v1.internal.routes.go
    │   │   └── router.go
    │   └── svc
    │       ├── repo.go
    │       └── svc.eventlogger.go
    ├── constants
    │   ├── env.constants.go
    │   └── utils.go
    ├── custom_parsers.conf
    ├── docker-compose.yml
    ├── dockerfile
    ├── fluent-bit
    │   ├── custom_parsers.conf
    │   ├── dockerfile
    │   └── fluent-bit.conf
    ├── fluent-bit.conf
    ├── go.mod
    ├── go.sum
    ├── logger
    │   ├── constants.go
    │   ├── eventlogger.go
    │   └── logger.go
    ├── logrotate.conf
    ├── logs
    │   ├── access.log
    │   ├── api_err.log
    │   ├── api_info.log
    │   ├── error.log
    │   └── events.log
    ├── main.go
    ├── middlewares
    │   ├── apilogger.go
    │   ├── auth.go
    │   └── ratelimiter.go
    ├── models
    │   ├── model.application
    │   │   └── v1.model.application.go
    │   └── model.eventlogger
    │       └── v1.eventlogger.go
    ├── test
    │   ├── dockerfile
    │   └── main.py
    └── utils
        ├── http.client
        │   └── client.go
        ├── response
        │   └── response.error.go
        └── utils.go

```