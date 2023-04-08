# Materialized Views and Koa.js

In this simple example, we will create a project using Kafka and the Koa framework in Node.js, which will demonstrate the purpose of Materialized Views. We will build a basic online voting system where users can vote for candidates, and we will keep real-time counts of the votes for each candidate. We will use the kafkajs library and node-kafka-streams to implement Kafka producers, consumers, and stream processing.

## Usage

```bash
curl -X POST -H "Content-Type: application/json" -d '{"candidate": "Jerry Shi"}' http://localhost:3000/vote
```
