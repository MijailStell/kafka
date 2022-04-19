# POC Kafka and Golang

## Setup

Run Kafka locally in docker

```bash
docker-compose up -d
```
![1](https://user-images.githubusercontent.com/1031887/163912481-6ef575a9-7503-471b-a622-13b26282dbc6.PNG)

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

![1](https://user-images.githubusercontent.com/1031887/163912255-0b8d5c63-7d67-4377-bf53-61b4bb3e8598.PNG)
