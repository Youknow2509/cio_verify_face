import { Injectable, OnModuleInit, OnModuleDestroy } from '@nestjs/common';

export interface ServiceInfo {
  serviceName: string;
  port: number;
  host: string;
  version: string;
  environment: string;
  status: 'UP' | 'DOWN';
  lastHeartbeat: Date;
}

@Injectable()
export class ServiceDiscoveryService implements OnModuleInit, OnModuleDestroy {
  private serviceInfo: ServiceInfo;
  private heartbeatInterval: NodeJS.Timeout;

  constructor(config: {
    serviceName: string;
    port: number;
    host: string;
    version: string;
    environment: string;
  }) {
    this.serviceInfo = {
      ...config,
      status: 'UP',
      lastHeartbeat: new Date(),
    };
  }

  onModuleInit() {
    this.startHeartbeat();
  }

  onModuleDestroy() {
    this.stopHeartbeat();
  }

  private startHeartbeat() {
    // Send heartbeat every 30 seconds
    this.heartbeatInterval = setInterval(() => {
      this.serviceInfo.lastHeartbeat = new Date();
      this.serviceInfo.status = 'UP';
      console.log(`[${new Date().toISOString()}] Identity Service heartbeat - Status: ${this.serviceInfo.status}`);
    }, 30000);
  }

  private stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
    }
    this.serviceInfo.status = 'DOWN';
  }

  getServiceInfo(): ServiceInfo {
    return { ...this.serviceInfo };
  }

  getServiceUrl(): string {
    return `http://${this.serviceInfo.host}:${this.serviceInfo.port}`;
  }

  getApiBaseUrl(): string {
    return `${this.getServiceUrl()}/api/v1`;
  }

  // Method to register with service registry (implement based on your service registry)
  async registerWithRegistry() {
    // Implementation depends on your service registry (Consul, Eureka, etc.)
    console.log(`Registering service: ${this.serviceInfo.serviceName} at ${this.getServiceUrl()}`);
  }

  // Method to deregister from service registry
  async deregisterFromRegistry() {
    // Implementation depends on your service registry
    console.log(`Deregistering service: ${this.serviceInfo.serviceName}`);
  }
}
