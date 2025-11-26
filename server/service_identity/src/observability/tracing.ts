import { NodeSDK } from '@opentelemetry/sdk-node';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { Resource } from '@opentelemetry/resources';
import {
    SEMRESATTRS_SERVICE_NAME,
    SEMRESATTRS_SERVICE_VERSION,
} from '@opentelemetry/semantic-conventions';

let sdk: NodeSDK | null = null;

export const initTracing = (
    serviceName: string,
    otlpEndpoint: string,
    environment: string = 'development'
): NodeSDK | null => {
    try {
        const traceExporter = new OTLPTraceExporter({
            url: otlpEndpoint,
        });

        sdk = new NodeSDK({
            resource: new Resource({
                [SEMRESATTRS_SERVICE_NAME]: serviceName,
                [SEMRESATTRS_SERVICE_VERSION]: '1.0.0',
                environment,
            }),
            traceExporter,
            instrumentations: [
                getNodeAutoInstrumentations({
                    '@opentelemetry/instrumentation-fs': {
                        enabled: false,
                    },
                }),
            ],
        });

        sdk.start();
        console.log('OpenTelemetry tracing initialized successfully');
        return sdk;
    } catch (error) {
        console.error('Failed to initialize tracing:', error);
        return null;
    }
};

export const shutdownTracing = async (): Promise<void> => {
    if (sdk) {
        try {
            await sdk.shutdown();
            console.log('OpenTelemetry tracing shut down successfully');
        } catch (error) {
            console.error('Error shutting down tracing:', error);
        }
    }
};

// Graceful shutdown
process.on('SIGTERM', async () => {
    await shutdownTracing();
    process.exit(0);
});
