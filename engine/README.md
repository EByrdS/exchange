# Engine

Kafka:

```
kafka-topics --create \
  --topic engine.DOLS.MEEM \
  --partitions 1 \
  --replication-factor 1 \
  --bootstrap-server localhost:9092
```

```
kafka-topics --create \
  --topic engine.DOLS.MEEM.volumes \
  --partitions 1 \
  --replication-factor 1 \
  --bootstrap-server localhost:9092
```

```
kafka-topics --create \
  --topic engine.DOLS.MEEM.orders \
  --partitions 1 \
  --replication-factor 1 \
  --bootstrap-server localhost:9092
```

```
kafka-topics --create \
  --topic engine.DOLS.MEEM.matches \
  --partitions 1 \
  --replication-factor 1 \
  --bootstrap-server localhost:9092
```
