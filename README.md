# POC Kafka and Golang

## 1. Setup for Local

[Descargar](https://www.apache.org/dyn/closer.cgi?path=/kafka/3.1.0/kafka_2.13-3.1.0.tgz) la última versión de Kafka y extraerla

```bash
$ tar -xzf kafka_2.13-3.1.0.tgz
$ cd kafka_2.13-3.1.0
```

Ejecutar los siguientes comandos para iniciar todos los servicios en el orden correcto:
```bash
$ bin/zookeeper-server-start.sh config/zookeeper.properties
```
![1](https://user-images.githubusercontent.com/1031887/163938848-d92fcba5-fe49-451f-a612-0d024c872cd5.PNG)

Abrir otra sesión de terminal y ejecutar el siguiente comando:
```bash
$ bin/kafka-server-start.sh config/server.properties
```
![2](https://user-images.githubusercontent.com/1031887/163939141-06bd9699-40b2-47a7-99af-eec7eb356bec.PNG)

Modificar los archivos consumer/config.yaml y producer/config.yaml.

```bash
LOCAL: "1"
```

Levantar producer y consumer

![1](https://user-images.githubusercontent.com/1031887/163939895-0e2ada26-dceb-4490-b019-92c3edac13bd.PNG)

Así debe mostrarse la consola del producer.

![3](https://user-images.githubusercontent.com/1031887/163939950-25bfa0ba-070e-466c-8bf3-e028c978b71a.PNG)

Ahora podemos realizar una petición al api del producer:

```bash
POST http://localhost:8080
Content-Type: application/json

{
    "data":"Hello-api1"
}
```

![2](https://user-images.githubusercontent.com/1031887/163940008-723e0e9a-004e-4235-8c17-7ada1e78500e.PNG)

Y podemos ver que el consumer muestra el mensaje leído en la consola.

![4](https://user-images.githubusercontent.com/1031887/163940044-bb58d903-0892-45d5-8df8-0f427d9595df.PNG)




## 2. Setup for Docker

Modificar los archivos consumer/config.yaml y producer/config.yaml.

```bash
LOCAL: "0"
```

Ejecutar kafka en docker localmente

```bash
docker-compose up -d
```
![1](https://user-images.githubusercontent.com/1031887/163912481-6ef575a9-7503-471b-a622-13b26282dbc6.PNG)

### Others commands

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
