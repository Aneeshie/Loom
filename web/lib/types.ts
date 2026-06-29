export interface LogEntry {
  id: string;
  timestamp: string;
  level: 'info' | 'warn' | 'error';
  message: string;
  service: string;
  host: string;
  metadata?: Record<string, any>;
}
