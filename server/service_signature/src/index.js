const express = require('express');
const cors = require('cors');
const swaggerUi = require('swagger-ui-express');
const swaggerSpec = require('./config/swagger');
const { initTracing } = require('./observability/tracing');
const {
    metricsMiddleware,
    setupMetricsServer,
} = require('./observability/metrics');
require('dotenv').config();

// Initialize tracing if enabled
const observabilityEnabled = process.env.OBSERVABILITY_ENABLED === 'true';
const tracingEnabled = process.env.TRACING_ENABLED === 'true';
if (observabilityEnabled && tracingEnabled) {
    const serviceName = process.env.SERVICE_NAME || 'service_signature';
    const otlpEndpoint =
        process.env.OTLP_ENDPOINT || 'http://jaeger:4318/v1/traces';
    const environment = process.env.NODE_ENV || 'development';
    initTracing(serviceName, otlpEndpoint, environment);
}

const app = express();
const PORT = process.env.PORT || 3001;

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.static('uploads'));

// Metrics middleware (should be before routes)
if (observabilityEnabled) {
    app.use(metricsMiddleware);
}

// Swagger UI
app.use('/api-docs', swaggerUi.serve);
app.get(
    '/api-docs',
    swaggerUi.setup(swaggerSpec, {
        swaggerOptions: {
            urls: [
                {
                    url: '/swagger.json',
                    name: 'Signature Service API',
                },
            ],
        },
    })
);

// Swagger JSON endpoint
app.get('/swagger.json', (req, res) => {
    res.setHeader('Content-Type', 'application/json');
    res.send(swaggerSpec);
});

// Routes
const signatureRoutes = require('./routes/signatures');
app.use('/api/v1/signatures', signatureRoutes);

// Health check
app.get('/health', (req, res) => {
    res.json({ status: 'OK' });
});

// Error handler
app.use((err, req, res, next) => {
    console.error(err);
    res.status(500).json({
        message: 'Internal Server Error',
        error: err.message,
    });
});

// 404 handler
app.use((req, res) => {
    res.status(404).json({ message: 'Route not found' });
});

// Start metrics server if enabled
if (observabilityEnabled) {
    const metricsPort = parseInt(process.env.METRICS_PORT || '9090', 10);
    const metricsPath = process.env.METRICS_PATH || '/metrics';
    setupMetricsServer(metricsPort, metricsPath);
}

app.listen(PORT, () => {
    console.log(`Signature Service running on port ${PORT}`);
    console.log(`Swagger UI available at http://localhost:${PORT}/api-docs`);
    if (observabilityEnabled) {
        console.log(
            `Metrics: http://localhost:${process.env.METRICS_PORT || 9090}${
                process.env.METRICS_PATH || '/metrics'
            }`
        );
        if (tracingEnabled) {
            console.log(
                `Tracing: Enabled (OTLP endpoint: ${process.env.OTLP_ENDPOINT})`
            );
        }
    }
});
