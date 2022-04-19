# POC Kafka and Golang

## Setup

Run Kafka locally in docker

```bash
docker-compose up -d
```

```bash
/opt/kafka/bin/kafka-topics.sh --list --bootstrap-server localhost:9092
```

```bash
/opt/kafka/bin/kafka-console-consumer.sh --topic topic1 --from-beginning --bootstrap-server localhost:9092
```

```bash
/opt/kafka/bin/kafka-console-consumer.sh --topic topic1 --bootstrap-server localhost:9092
```

## Send Data JSON
```bash
POST http://localhost:8080
Content-Type: application/json

{
    "data":"Hello-api"
}
```