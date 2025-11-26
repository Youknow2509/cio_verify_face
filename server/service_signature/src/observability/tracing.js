const { NodeSDK } = require('@opentelemetry/sdk-node');
const {
    getNodeAutoInstrumentations,
} = require('@opentelemetry/auto-instrumentations-node');
const {
    OTLPTraceExporter,
} = require('@opentelemetry/exporter-trace-otlp-http');
const { Resource } = require('@opentelemetry/resources');
const {
    SEMRESATTRS_SERVICE_NAME,
    SEMRESATTRS_SERVICE_VERSION,
} = require('@opentelemetry/semantic-conventions');

let sdk = null;

const initTracing = (
    serviceName,
    otlpEndpoint,
    environment = 'development'
) => {
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

const shutdownTracing = async () => {
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

module.exports = {
    initTracing,
    shutdownTracing,
};
