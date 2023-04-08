const { Kafka } = require('kafkajs');

const kafka = new Kafka({
  clientId: 'vote-counter',
  brokers: ['localhost:9092'],
});

const consumer = kafka.consumer({ groupId: 'vote-counter-group' });
const producer = kafka.producer();
const inputTopic = 'vote-events';
const outputTopic = 'vote-results';

// const voteCounts: { [key: string]: number } = {};
// 使用 voteCounts 对象作为 Materialized View 来存储实时计算的候选人得票数
const voteCounts = {};

(async () => {
  await consumer.connect();
  await producer.connect();
  await consumer.subscribe({ topic: inputTopic });

  await consumer.run({
    eachMessage: async ({ message }) => {
      console.log('view-counter message:', message);
      const { candidate } = JSON.parse(message.value.toString());
      if (!voteCounts[candidate]) {
        voteCounts[candidate] = 0;
      }
      voteCounts[candidate] += 1;

      await producer.send({
        topic: outputTopic,
        messages: [
          { key: candidate, value: JSON.stringify(voteCounts[candidate]) },
        ],
      });
    },
  });
})();
