import express, { Express } from 'express';
import cors from 'cors';
import helmet from 'helmet';
import dotenv from 'dotenv';
import swaggerUi from 'swagger-ui-express';
import { swaggerSpec } from './config/swagger';
import apiRoutes from './routes';
import { errorHandler, notFoundHandler } from './middleware/errorHandler';
import { initTracing } from './observability/tracing';
import { metricsMiddleware, setupMetricsServer } from './observability/metrics';

dotenv.config();

// Initialize tracing if enabled
const observabilityEnabled = process.env.OBSERVABILITY_ENABLED === 'true';
const tracingEnabled = process.env.TRACING_ENABLED === 'true';
if (observabilityEnabled && tracingEnabled) {
    const serviceName = process.env.SERVICE_NAME || 'service_identity';
    const otlpEndpoint =
        process.env.OTLP_ENDPOINT || 'http://jaeger:4318/v1/traces';
    const environment = process.env.NODE_ENV || 'development';
    initTracing(serviceName, otlpEndpoint, environment);
}

const app: Express = express();
const PORT = process.env.PORT || 3001;

// Middleware
app.use(helmet());
app.use(
    cors({
        origin: (_origin, cb) => cb(null, true),
        credentials: true,
        methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS'],
        allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With'],
    })
);
app.options('*', cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Metrics middleware (should be before routes)
if (observabilityEnabled) {
    app.use(metricsMiddleware);
}

// Swagger documentation
app.use('/api-docs', swaggerUi.serve);
app.get('/api-docs', swaggerUi.setup(swaggerSpec, { explorer: true }));

// Health check endpoint
app.get('/health', (req, res) => {
    res.status(200).json({ status: 'OK', timestamp: new Date().toISOString() });
});

// API routes
app.use('/api/v1', apiRoutes);

// 404 handler
app.use(notFoundHandler);

// Error handler (must be last)
app.use(errorHandler);

// Start metrics server if enabled
if (observabilityEnabled) {
    const metricsPort = parseInt(process.env.METRICS_PORT || '9090', 10);
    const metricsPath = process.env.METRICS_PATH || '/metrics';
    setupMetricsServer(metricsPort, metricsPath);
}

// Start server
app.listen(PORT, () => {
    console.log(`Identity & Organization Service running on port ${PORT}`);
    console.log(`Environment: ${process.env.NODE_ENV || 'development'}`);
    console.log(`API Documentation: http://localhost:${PORT}/api-docs`);
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

export default app;
