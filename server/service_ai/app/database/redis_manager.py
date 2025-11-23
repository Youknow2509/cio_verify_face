"""
    Redis distributed cache manager module
"""

import redis
from redis.sentinel import Sentinel
from redis.cluster import RedisCluster
from typing import Optional
import logging
from app.core.config import settings

logger = logging.getLogger(__name__)


class RedisManager:
    def __init__(self):
        """Initialize Redis manager"""
        self._client: Optional[redis.Redis] = None
        self._connect()
    
    def check_connection(self) -> bool:
        """Check if Redis connection is alive"""
        try:
            return self.client.ping()
        except Exception as e:
            logger.error(f"Redis connection check failed: {e}")
            return False
    
    def _connect(self):
        """Connect to Redis server based on configuration"""
        try:
            # Type 1: Standalone Redis
            if settings.REDIS_TYPE == 1:
                self._client = self._connect_standalone()
                logger.info(f"Connected to standalone Redis at {settings.REDIS_HOST}:{settings.REDIS_PORT}")
            
            # Type 2: Sentinel Redis
            elif settings.REDIS_TYPE == 2:
                self._client = self._connect_sentinel()
                logger.info(f"Connected to Redis Sentinel with master: {settings.REDIS_MASTER_NAME}")
            
            # Type 3: Cluster Redis
            elif settings.REDIS_TYPE == 3:
                self._client = self._connect_cluster()
                logger.info(f"Connected to Redis Cluster")
            
            else:
                raise ValueError(f"Invalid REDIS_TYPE: {settings.REDIS_TYPE}. Must be 1, 2, or 3.")
            
            # Test connection
            self._client.ping()
            logger.info("Redis connection successful")
            
        except Exception as e:
            logger.error(f"Failed to connect to Redis: {e}")
            raise
    
    def _connect_standalone(self) -> redis.Redis:
        """Connect to standalone Redis server"""
        connection_params = {
            'host': settings.REDIS_HOST,
            'port': settings.REDIS_PORT,
            'db': settings.REDIS_DB,
            'password': settings.REDIS_PASSWORD if settings.REDIS_PASSWORD else None,
            'decode_responses': True,
            'max_connections': settings.REDIS_POOL_SIZE,
            'socket_keepalive': True,
            'socket_connect_timeout': 5,
            'socket_timeout': 5,
            'retry_on_timeout': True,
            'health_check_interval': 30,
        }
        
        # Add TLS configuration if enabled
        if settings.REDIS_USE_TLS:
            connection_params.update({
                'ssl': True,
                'ssl_certfile': settings.REDIS_CERT_PATH,
                'ssl_keyfile': settings.REDIS_KEY_PATH,
            })
        
        return redis.Redis(**connection_params)
    
    def _connect_sentinel(self) -> redis.Redis:
        """Connect to Redis Sentinel"""
        # Parse sentinel addresses
        sentinel_list = []
        for addr in settings.REDIS_SENTINEL_ADDRS:
            host, port = addr.split(':')
            sentinel_list.append((host, int(port)))
        
        sentinel_params = {
            'password': settings.REDIS_PASSWORD if settings.REDIS_PASSWORD else None,
            'socket_keepalive': True,
            'socket_connect_timeout': 5,
            'socket_timeout': 5,
        }
        
        # Add TLS configuration if enabled
        if settings.REDIS_USE_TLS:
            sentinel_params.update({
                'ssl': True,
                'ssl_certfile': settings.REDIS_CERT_PATH,
                'ssl_keyfile': settings.REDIS_KEY_PATH,
            })
        
        sentinel = Sentinel(
            sentinel_list,
            sentinel_kwargs=sentinel_params,
            socket_timeout=5,
        )
        
        # Get master connection
        master = sentinel.master_for(
            settings.REDIS_MASTER_NAME,
            db=settings.REDIS_DB,
            decode_responses=True,
            max_connections=settings.REDIS_POOL_SIZE,
        )
        
        return master
    
    def _connect_cluster(self) -> RedisCluster:
        """Connect to Redis Cluster"""
        # Parse cluster addresses
        startup_nodes = []
        for addr in settings.REDIS_CLUSTER_ADDRS:
            host, port = addr.split(':')
            startup_nodes.append({'host': host, 'port': int(port)})
        
        cluster_params = {
            'startup_nodes': startup_nodes,
            'password': settings.REDIS_PASSWORD if settings.REDIS_PASSWORD else None,
            'decode_responses': True,
            'max_connections': settings.REDIS_POOL_SIZE,
            'max_connections_per_node': True,
            'read_from_replicas': settings.REDIS_ROUTE_BY_LATENCY,
            'skip_full_coverage_check': True,
            'socket_keepalive': True,
            'socket_connect_timeout': 5,
            'socket_timeout': 5,
            'retry_on_timeout': True,
        }
        
        # Add TLS configuration if enabled
        if settings.REDIS_USE_TLS:
            cluster_params.update({
                'ssl': True,
                'ssl_certfile': settings.REDIS_CERT_PATH,
                'ssl_keyfile': settings.REDIS_KEY_PATH,
            })
        
        return RedisCluster(**cluster_params)
    
    @property
    def client(self) -> redis.Redis:
        """Get Redis client instance"""
        if self._client is None:
            self._connect()
        return self._client
    
    def get(self, key: str) -> Optional[str]:
        """Get value by key"""
        try:
            return self.client.get(key)
        except Exception as e:
            logger.error(f"Redis GET error for key '{key}': {e}")
            raise
    
    def set(self, key: str, value: str, ex: Optional[int] = None) -> bool:
        """Set key-value pair with optional expiration"""
        try:
            return self.client.set(key, value, ex=ex)
        except Exception as e:
            logger.error(f"Redis SET error for key '{key}': {e}")
            raise
    
    def delete(self, *keys: str) -> int:
        """Delete one or more keys"""
        try:
            return self.client.delete(*keys)
        except Exception as e:
            logger.error(f"Redis DELETE error: {e}")
            raise
    
    def exists(self, *keys: str) -> int:
        """Check if keys exist"""
        try:
            return self.client.exists(*keys)
        except Exception as e:
            logger.error(f"Redis EXISTS error: {e}")
            raise
    
    def expire(self, key: str, seconds: int) -> bool:
        """Set expiration on key"""
        try:
            return self.client.expire(key, seconds)
        except Exception as e:
            logger.error(f"Redis EXPIRE error for key '{key}': {e}")
            raise
    
    def close(self):
        """Close Redis connection"""
        if self._client:
            try:
                self._client.close()
                logger.info("Redis connection closed")
            except Exception as e:
                logger.error(f"Error closing Redis connection: {e}")
    
    def __enter__(self):
        """Context manager entry"""
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit"""
        self.close()