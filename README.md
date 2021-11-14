# k2rt - Kafka To Redis-Timeseries

Dump kafka messages to Redis-Timeseries by pattern

## Setup

1. Get k2rt (`ghcr.io/slntopp/k2rt` Docker image recommended)
2. Set environment variables: `KAFKA_HOST`, `REDIS_HOST`, `TOPIC`
3. Start

Docker Compose service example:

```yml
  k2rt:
    image: ghcr.io/slntopp/k2rt
    environment:
      KAFKA_HOST: kafka:29092
      REDIS_HOST: timeseries:6379
      TOPIC: shadow.reported-state.delta
```

## What's going to happen

`k2rt` will read all messages from `TOPIC` as JSON and save it to Redis Timeseries as:

Key    = `{Event Key}:{Data JSON Key}:`
Value  = Data from JSON(only int, float and bool values are saved)
Labels =
  Label Key   = `Data JSON Key`
  Label Value = Data from JSON(string values)
