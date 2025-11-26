const { Registry, Counter, Histogram, Gauge } = require('prom-client');
const express = require('express');

// Create a Registry to register the metrics
const register = new Registry();

// HTTP metrics
const httpRequestsTotal = new Counter({
    name: 'http_requests_total',
    help: 'Total number of HTTP requests',
    labelNames: ['method', 'route', 'status_code'],
    registers: [register],
});

const httpRequestDuration = new Histogram({
    name: 'http_request_duration_seconds',
    help: 'Duration of HTTP requests in seconds',
    labelNames: ['method', 'route', 'status_code'],
    buckets: [0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5],
    registers: [register],
});

const httpRequestSize = new Histogram({
    name: 'http_request_size_bytes',
    help: 'Size of HTTP requests in bytes',
    labelNames: ['method', 'route'],
    buckets: [100, 1000, 5000, 10000, 50000, 100000, 500000, 1000000],
    registers: [register],
});

const httpResponseSize = new Histogram({
    name: 'http_response_size_bytes',
    help: 'Size of HTTP responses in bytes',
    labelNames: ['method', 'route'],
    buckets: [100, 1000, 5000, 10000, 50000, 100000, 500000, 1000000],
    registers: [register],
});

// Application metrics
const activeConnections = new Gauge({
    name: 'active_connections',
    help: 'Number of active connections',
    registers: [register],
});

// Middleware to collect metrics
const metricsMiddleware = (req, res, next) => {
    const start = Date.now();

    // Track request size
    const requestSize = parseInt(req.get('content-length') || '0', 10);
    if (requestSize > 0) {
        httpRequestSize.observe(
            { method: req.method, route: req.route?.path || req.path },
            requestSize
        );
    }

    // Track active connections
    activeConnections.inc();

    // Override res.end to capture response
    const originalEnd = res.end;
    res.end = function (...args) {
        const duration = (Date.now() - start) / 1000;
        const route = req.route?.path || req.path;

        // Record metrics
        httpRequestsTotal.inc({
            method: req.method,
            route,
            status_code: res.statusCode,
        });

        httpRequestDuration.observe(
            { method: req.method, route, status_code: res.statusCode },
            duration
        );

        // Track response size
        const responseSize = parseInt(res.get('content-length') || '0', 10);
        if (responseSize > 0) {
            httpResponseSize.observe(
                { method: req.method, route },
                responseSize
            );
        }

        activeConnections.dec();

        return originalEnd.apply(this, args);
    };

    next();
};

// Metrics endpoint handler
const metricsHandler = async (req, res) => {
    res.set('Content-Type', register.contentType);
    res.end(await register.metrics());
};

// Setup metrics server
const setupMetricsServer = (port, path = '/metrics') => {
    const metricsApp = express();

    metricsApp.get(path, metricsHandler);
    metricsApp.get('/health', (req, res) => {
        res.status(200).json({ status: 'OK' });
    });

    metricsApp.listen(port, () => {
        console.log(`Metrics server listening on port ${port}${path}`);
    });
};

module.exports = {
    register,
    metricsMiddleware,
    metricsHandler,
    setupMetricsServer,
};
