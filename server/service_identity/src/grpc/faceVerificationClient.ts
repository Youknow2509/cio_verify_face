import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import path from 'path';

// Wrapper client for FaceVerification service
const PROTO_PATH = path.join(
    __dirname,
    '..',
    'proto',
    'face_verification.proto'
);

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const faceProto: any =
    grpc.loadPackageDefinition(packageDefinition).face_verification;

export interface EnrollParams {
    user_id: string;
    company_id: string;
    image_data: Buffer;
    make_primary?: boolean;
    metadata?: Record<string, string>;
}

export interface EnrollResult {
    status: string;
    profile_id?: string;
    message?: string;
    quality_score?: number;
    embedding?: number[];
    embedding_version?: string;
}

class FaceVerificationClient {
    private client: grpc.Client;

    constructor(address: string) {
        this.client = new faceProto.FaceVerification(
            address,
            grpc.credentials.createInsecure()
        );
    }

    enrollFace(params: EnrollParams): Promise<EnrollResult> {
        return new Promise((resolve, reject) => {
            // Method names are exposed in lowerCamelCase by proto-loader
            (this.client as any).enrollFace(
                {
                    user_id: params.user_id,
                    company_id: params.company_id,
                    image_data: params.image_data,
                    make_primary: params.make_primary || false,
                    metadata: params.metadata || {},
                },
                (err: grpc.ServiceError, response: EnrollResult) => {
                    if (err) return reject(err);
                    resolve(response);
                }
            );
        });
    }
}

// Singleton instance; address should be configured via env FACE_AI_GRPC_ADDR (e.g. "localhost:50051")
const grpcAddress = process.env.FACE_AI_GRPC_ADDR || 'localhost:50051';
export const faceVerificationClient = new FaceVerificationClient(grpcAddress);
