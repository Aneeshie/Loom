import { LogEntry } from './types';

export const getMockLogs = (): LogEntry[] => {
  const now = new Date();
  
  return [
    {
      id: '1',
      timestamp: new Date(now.getTime() - 2 * 60 * 1000).toISOString(), // 2 minutes ago
      level: 'error',
      message: 'Connection timeout to database pool [postgres_db]',
      service: 'payment-gateway',
      host: 'prod-ap-south-1a',
      metadata: {
        request_id: 'req_8f3a9d2c',
        http_method: 'POST',
        http_path: '/api/v1/payments/charge',
        http_status: 504,
        duration_ms: 5000,
        db_queries: 1,
        ip_address: '192.168.1.42',
        error_code: 'ERR_DB_TIMEOUT',
        exception: 'ConnectionPoolTimeoutException: Timeout waiting for idle connection in pool after 5000ms'
      }
    },
    {
      id: '2',
      timestamp: new Date(now.getTime() - 4 * 60 * 1000).toISOString(), // 4 minutes ago
      level: 'info',
      message: 'User session validated and authentication successful',
      service: 'auth-service',
      host: 'prod-ap-south-1b',
      metadata: {
        request_id: 'req_a2b3c4d5',
        http_method: 'GET',
        http_path: '/api/v1/auth/session',
        http_status: 200,
        duration_ms: 45,
        user_id: 'usr_99218',
        ip_address: '203.0.113.88'
      }
    },
    {
      id: '3',
      timestamp: new Date(now.getTime() - 8 * 60 * 1000).toISOString(), // 8 minutes ago
      level: 'warn',
      message: 'High memory usage (88%) detected on node',
      service: 'worker-node-1',
      host: 'prod-ap-south-1a',
      metadata: {
        cpu_usage_pct: 74.2,
        memory_used_bytes: 15110992384,
        memory_total_bytes: 17179869184,
        active_threads: 242,
        disk_io_ops: 890
      }
    },
    {
      id: '4',
      timestamp: new Date(now.getTime() - 12 * 60 * 1000).toISOString(), // 12 minutes ago
      level: 'error',
      message: 'Upstream payment gateway API timeout after 5000ms',
      service: 'payment-gateway',
      host: 'prod-ap-south-1a',
      metadata: {
        request_id: 'req_1a2b3c4d',
        http_method: 'POST',
        http_path: '/api/v1/payments/refund',
        http_status: 502,
        duration_ms: 5012,
        provider: 'Stripe',
        error_code: 'GATEWAY_TIMEOUT'
      }
    },
    {
      id: '5',
      timestamp: new Date(now.getTime() - 25 * 60 * 1000).toISOString(), // 25 minutes ago
      level: 'info',
      message: 'Cache successfully invalidated for user profile stats',
      service: 'cache-manager',
      host: 'prod-ap-south-1b',
      metadata: {
        keys_invalidated: 42,
        cache_driver: 'Redis',
        operation: 'flush_pattern',
        pattern: 'users:stats:*',
        duration_ms: 12
      }
    },
    {
      id: '6',
      timestamp: new Date(now.getTime() - 40 * 60 * 1000).toISOString(), // 40 minutes ago
      level: 'warn',
      message: 'Slow query detected on users table (duration: 1.2s)',
      service: 'user-db-pool',
      host: 'prod-ap-south-1a',
      metadata: {
        query: 'SELECT * FROM users WHERE status = \'active\' ORDER BY last_login_at DESC LIMIT 100 OFFSET 5000',
        duration_ms: 1205,
        rows_returned: 100,
        lock_wait_ms: 4
      }
    },
    {
      id: '7',
      timestamp: new Date(now.getTime() - 75 * 60 * 1000).toISOString(), // 1.25 hours ago
      level: 'info',
      message: 'Received webhook notification for invoice subscription.paid',
      service: 'payment-gateway',
      host: 'prod-ap-south-1b',
      metadata: {
        request_id: 'req_f9a8b7c6',
        event_type: 'invoice.subscription.paid',
        customer_id: 'cus_Hj3k8F91s',
        amount_cents: 2900,
        currency: 'usd'
      }
    },
    {
      id: '8',
      timestamp: new Date(now.getTime() - 110 * 60 * 1000).toISOString(), // ~2 hours ago
      level: 'error',
      message: 'Failed to generate PDF invoice: out of memory disk space',
      service: 'analytics-api',
      host: 'staging-eu-west-1a',
      metadata: {
        invoice_id: 'inv_10284',
        disk_free_bytes: 0,
        temp_directory: '/tmp/pdf_worker/',
        error_code: 'ENOSPC'
      }
    },
    {
      id: '9',
      timestamp: new Date(now.getTime() - 180 * 60 * 1000).toISOString(), // 3 hours ago
      level: 'info',
      message: 'Daily system health check completed with 0 critical issues',
      service: 'worker-node-1',
      host: 'staging-eu-west-1a',
      metadata: {
        tasks_run: 15,
        successful_tasks: 15,
        failed_tasks: 0,
        uptime_seconds: 864000
      }
    },
    {
      id: '10',
      timestamp: new Date(now.getTime() - 320 * 60 * 1000).toISOString(), // ~5 hours ago
      level: 'warn',
      message: 'Failed login attempt from unrecognized IP address range',
      service: 'auth-service',
      host: 'prod-ap-south-1b',
      metadata: {
        username: 'admin_test',
        ip_address: '198.51.100.12',
        attempt_count: 3,
        location: 'Unknown (VPN)'
      }
    },
    {
      id: '11',
      timestamp: new Date(now.getTime() - 720 * 60 * 1000).toISOString(), // 12 hours ago
      level: 'info',
      message: 'Analytics report exported successfully to Amazon S3',
      service: 'analytics-api',
      host: 'prod-ap-south-1a',
      metadata: {
        report_type: 'monthly_revenue_summary',
        s3_bucket: 'loom-billing-reports',
        s3_key: '2026/06/revenue_summary.csv',
        file_size_bytes: 4882109,
        duration_ms: 12890
      }
    },
    {
      id: '12',
      timestamp: new Date(now.getTime() - 1500 * 60 * 1000).toISOString(), // 25 hours ago (~1 day)
      level: 'error',
      message: 'Failed to verify JWT signature: token has expired',
      service: 'auth-service',
      host: 'prod-ap-south-1b',
      metadata: {
        request_id: 'req_jwt_err_9a',
        http_method: 'GET',
        http_path: '/api/v1/user/settings',
        token_issuer: 'auth-service',
        expiry_timestamp: '2026-06-28T08:00:00Z'
      }
    },
    {
      id: '13',
      timestamp: new Date(now.getTime() - 2800 * 60 * 1000).toISOString(), // ~2 days ago
      level: 'info',
      message: 'Background indexing scheduler started',
      service: 'worker-node-1',
      host: 'dev-us-east-1',
      metadata: {
        scheduler_name: 'es_reindex_worker',
        batch_size: 500,
        concurrency: 4
      }
    },
    {
      id: '14',
      timestamp: new Date(now.getTime() - 6000 * 60 * 1000).toISOString(), // ~4 days ago
      level: 'warn',
      message: 'Rate limit threshold reached for developer API key',
      service: 'analytics-api',
      host: 'dev-us-east-1',
      metadata: {
        api_key_id: 'key_dev_0x9b32',
        request_count_minute: 120,
        limit_allowed_minute: 100,
        ip_address: '45.79.112.5'
      }
    },
    {
      id: '15',
      timestamp: new Date(now.getTime() - 9000 * 60 * 1000).toISOString(), // ~6 days ago
      level: 'info',
      message: 'Database backup successfully uploaded to Glacier',
      service: 'user-db-pool',
      host: 'prod-ap-south-1a',
      metadata: {
        backup_id: 'backup_20260623',
        backup_size_bytes: 42949672960,
        compression_ratio: '3.4x',
        glacier_vault: 'loom-vault-db'
      }
    }
  ];
};
