// redis-counter.js

const redis = require('redis');
const { promisify } = require('util');
const amqp = require('amqplib');

const client = redis.createClient();
const incrAsync = promisify(client.incr).bind(client);
const lpushAsync = promisify(client.lpush).bind(client);
const lpopAsync = promisify(client.lpop).bind(client);
const setnxAsync = promisify(client.setnx).bind(client);
const publishAsync = promisify(client.publish).bind(client);

const QUEUE_NAME = 'update_queue';
const LOCK_KEY = 'counter_lock';

async function updateCounter() {
  try {
    const isLocked = await setnxAsync(LOCK_KEY, 'locked');
    
    if (isLocked === 0) {
      console.log('Failed to acquire lock. Another process is updating the counter.');
      return;
    }

    const counter = await incrAsync('counter');
    console.log(`Counter: ${counter}`);

    // Notify other processes about the successful update
    await publishAsync('counter_updated', 'Counter updated successfully');

    // Release the lock
    await client.del(LOCK_KEY);
  } catch (error) {
    console.error(`Error updating counter: ${error}`);
  }
}

async function processQueue() {
  while (true) {
    const message = await lpopAsync(QUEUE_NAME);
    if (!message) break;

    console.log(`Processing message: ${message}`);
    await updateCounter();
  }
}

async function main() {
  await client.flushall(); // Clear Redis data
  await client.set('counter', 0);

  const connection = await amqp.connect('amqp://localhost');
  const channel = await connection.createChannel();

  await channel.assertQueue(QUEUE_NAME);

  // Start two separate processes or threads
  setInterval(async () => {
    await lpushAsync(QUEUE_NAME, 'Increment from Process 1');
  }, 1000);

  setInterval(async () => {
    await lpushAsync(QUEUE_NAME, 'Increment from Process 2');
  }, 1200);

  // Start processing the queue
  processQueue();
}

main();
