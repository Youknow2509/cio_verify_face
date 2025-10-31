import swaggerJsdoc from 'swagger-jsdoc';

const options = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: 'Identity & Organization Service API',
      version: '1.0.0',
      description: 'API for managing companies, users, employees and face data',
      contact: {
        name: 'API Support',
        email: 'support@example.com',
      },
    },
    servers: [
      {
        url: 'http://localhost:3001',
        description: 'Development server',
      },
      {
        url: 'https://api.example.com',
        description: 'Production server',
      },
    ],
    components: {
      schemas: {
        Company: {
          type: 'object',
          properties: {
            company_id: {
              type: 'string',
              format: 'uuid',
              description: 'Unique company identifier',
            },
            name: {
              type: 'string',
              description: 'Company name',
            },
            address: {
              type: 'string',
              description: 'Company address',
            },
            phone: {
              type: 'string',
              description: 'Company phone number',
            },
            email: {
              type: 'string',
              format: 'email',
              description: 'Company email',
            },
            website: {
              type: 'string',
              format: 'uri',
              description: 'Company website',
            },
            status: {
              type: 'integer',
              enum: [0, 1, 2],
              description: '0: Inactive, 1: Active, 2: Suspended',
            },
            subscription_plan: {
              type: 'integer',
              enum: [0, 1, 2],
              description: '0: Basic, 1: Premium, 2: Enterprise',
            },
            subscription_start_date: {
              type: 'string',
              format: 'date',
            },
            subscription_end_date: {
              type: 'string',
              format: 'date',
            },
            max_employees: {
              type: 'integer',
              description: 'Maximum number of employees',
            },
            max_devices: {
              type: 'integer',
              description: 'Maximum number of devices',
            },
            created_at: {
              type: 'string',
              format: 'date-time',
            },
            updated_at: {
              type: 'string',
              format: 'date-time',
            },
          },
          required: ['company_id', 'name', 'status'],
        },
        User: {
          type: 'object',
          properties: {
            user_id: {
              type: 'string',
              format: 'uuid',
            },
            email: {
              type: 'string',
              format: 'email',
            },
            phone: {
              type: 'string',
            },
            full_name: {
              type: 'string',
            },
            avatar_url: {
              type: 'string',
              format: 'uri',
            },
            role: {
              type: 'integer',
              enum: [0, 1, 2],
              description: '0: SYSTEM_ADMIN, 1: COMPANY_ADMIN, 2: EMPLOYEE',
            },
            status: {
              type: 'integer',
              enum: [0, 1, 2],
              description: '0: ACTIVE, 1: INACTIVE, 2: SUSPENDED',
            },
            is_locked: {
              type: 'boolean',
            },
            created_at: {
              type: 'string',
              format: 'date-time',
            },
            updated_at: {
              type: 'string',
              format: 'date-time',
            },
          },
          required: ['user_id', 'email', 'phone', 'full_name', 'role'],
        },
        FaceData: {
          type: 'object',
          properties: {
            fid: {
              type: 'string',
              format: 'uuid',
            },
            user_id: {
              type: 'string',
              format: 'uuid',
            },
            image_url: {
              type: 'string',
              format: 'uri',
            },
            face_encoding: {
              type: 'string',
            },
            quality_score: {
              type: 'number',
              minimum: 0,
              maximum: 1,
            },
            created_at: {
              type: 'string',
              format: 'date-time',
            },
            updated_at: {
              type: 'string',
              format: 'date-time',
            },
          },
          required: ['fid', 'user_id', 'image_url'],
        },
        ApiResponse: {
          type: 'object',
          properties: {
            success: {
              type: 'boolean',
            },
            message: {
              type: 'string',
            },
            data: {
              type: 'object',
            },
            statusCode: {
              type: 'integer',
            },
          },
        },
        ErrorResponse: {
          type: 'object',
          properties: {
            success: {
              type: 'boolean',
              default: false,
            },
            message: {
              type: 'string',
            },
            error: {
              type: 'string',
            },
            statusCode: {
              type: 'integer',
            },
          },
        },
      },
    },
  },
  apis: ['./src/routes/*.ts'],
};

export const swaggerSpec = swaggerJsdoc(options);
