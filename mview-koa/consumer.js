const { Kafka } = require('kafkajs');
const Koa = require('koa');
const Router = require('koa-router');

const voteCounts = {};

const kafka = new Kafka({
  clientId: 'voting-consumer',
  brokers: ['localhost:9092'],
});

const consumer = kafka.consumer({ groupId: 'voting-consumer-group' });
const topic = 'vote-results';

const app = new Koa();
const router = new Router();

router.get('/results', async (ctx) => {
  ctx.status = 200;
  ctx.body = voteCounts;
});

app.use(router.routes());

app.listen(3001, async () => {
  console.log('Server is running on http://localhost:3001');

  await consumer.connect();
  await consumer.subscribe({ topic });

  consumer.run({
    eachMessage: async ({ message }) => {
      console.log('consumer message:', message);
      const key = message.key.toString();
      const value = message.value.toString();
      voteCounts[key] = parseInt(value, 10);
    },
  });
});
