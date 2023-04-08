const { Kafka } = require('kafkajs');
const Koa = require('koa');
const Router = require('koa-router');
const { koaBody } = require('koa-body');

const kafka = new Kafka({
  clientId: 'voting-producer',
  brokers: ['localhost:9092'],
});

const producer = kafka.producer();
const topic = 'vote-events';

const app = new Koa();
const router = new Router();

router.post('/vote', async (ctx) => {
  const { candidate } = ctx.request.body;

  if (!candidate) {
    ctx.status = 400;
    ctx.body = { error: 'Missing candidate' };
    return;
  }

  await producer.connect();
  await producer.send({
    topic,
    messages: [{ value: JSON.stringify({ candidate, timestamp: Date.now() }) }],
  });

  ctx.status = 200;
  ctx.body = { success: true };
});

app.use(koaBody());
app.use(router.routes());

app.listen(3000, () => {
  console.log('Server is running on http://localhost:3000');
});
