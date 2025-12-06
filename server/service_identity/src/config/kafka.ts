import { Kafka, logLevel, Producer } from 'kafkajs';

let kafkaProducer: Producer | null = null;

const KAFKA_USER_EVENTS_TOPIC =
    process.env.KAFKA_USER_EVENTS_TOPIC || 'user-events';
const KAFKA_PASSWORD_RESET_TOPIC =
    process.env.KAFKA_PASSWORD_RESET_TOPIC || 'user.password_reset';

export async function getKafkaProducer(): Promise<Producer> {
    if (kafkaProducer) {
        return kafkaProducer;
    }

    const kafka = new Kafka({
        clientId: 'service-identity',
        brokers: (process.env.KAFKA_BROKERS || 'localhost:9092').split(','),
        logLevel: logLevel.ERROR,
        retry: {
            initialRetryTime: 100,
            retries: 8,
        },
    });

    kafkaProducer = kafka.producer();
    await kafkaProducer.connect();

    // Handle graceful shutdown
    process.on('SIGTERM', async () => {
        if (kafkaProducer) {
            await kafkaProducer.disconnect();
        }
    });

    process.on('SIGINT', async () => {
        if (kafkaProducer) {
            await kafkaProducer.disconnect();
        }
    });

    return kafkaProducer;
}

export async function sendToKafka(
    topic: string,
    messages: any[]
): Promise<void> {
    try {
        const producer = await getKafkaProducer();
        await producer.send({
            topic: topic || KAFKA_USER_EVENTS_TOPIC,
            messages: messages.map((msg) => ({
                value: JSON.stringify(msg),
            })),
        });
        console.log(
            `Message sent to Kafka topic: ${topic || KAFKA_USER_EVENTS_TOPIC}`
        );
    } catch (error) {
        console.error(`Failed to send message to Kafka: ${error}`);
        throw error;
    }
}

export function getKafkaTopics() {
    return {
        userEvents: KAFKA_USER_EVENTS_TOPIC,
        passwordReset: KAFKA_PASSWORD_RESET_TOPIC,
    };
}
