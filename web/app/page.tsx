"use client";

import { useState, useMemo, useEffect } from 'react';
import SearchBar from '@/components/SearchBar';
import IntentCard from '@/components/IntentCard';
import LogList from '@/components/LogList';
import { LogEntry } from '@/lib/types';
import { Terminal, Activity, Cpu, Server, ShieldCheck, Loader2, AlertCircle } from 'lucide-react';

const SERVICES_LIST = [
  'payment-gateway',
  'auth-service',
  'worker-node-1',
  'cache-manager',
  'analytics-api',
  'user-db-pool',
];

const HOSTS_LIST = [
  'prod-ap-south-1a',
  'prod-ap-south-1b',
  'staging-eu-west-1a',
  'dev-us-east-1',
];

export default function Dashboard() {
  // Filter States
  const [searchQuery, setSearchQuery] = useState('');
  const [level, setLevel] = useState('all');
  const [service, setService] = useState('');
  const [host, setHost] = useState('');
  const [since, setSince] = useState('all');

  // API State
  const [activeLogs, setActiveLogs] = useState<LogEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLive, setIsLive] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  // Main fetch runner
  const fetchLiveLogs = async (queryText: string, lvl: string, srv: string, hst: string, duration: string) => {
    setIsLoading(true);
    setErrorMsg(null);
    try {
      const payload = {
        query: queryText,
        level: lvl === 'all' ? '' : lvl.toUpperCase(),
        service: srv,
        host: hst,
        since: duration === 'all' ? '' : duration,
      };

      const res = await fetch('/api/v1/query', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        throw new Error(`HTTP status ${res.status}`);
      }

      const data = await res.json();
      
      if (data && Array.isArray(data.logs)) {
        const apiLogs: LogEntry[] = data.logs.map((log: any) => ({
          id: String(log.id),
          timestamp: new Date(log.timestamp * 1000).toISOString(),
          level: (log.level || 'info').toLowerCase() as 'info' | 'warn' | 'error',
          message: log.message || '',
          service: log.service_name || '',
          host: log.host || '',
          metadata: {
            db_id: log.id,
            service_name: log.service_name,
            host: log.host,
            timestamp: log.timestamp,
            level: log.level,
            intent_extracted: data.intent
          }
        }));
        
        setActiveLogs(apiLogs);
        setIsLive(true);
      } else {
        throw new Error("Missing logs field");
      }
    } catch (err) {
      console.error("Go server API request failed:", err);
      setIsLive(false);
      setErrorMsg("Unable to connect to the Go Observability API server. Please check that your backend is running.");
      setActiveLogs([]);
    } finally {
      setIsLoading(false);
    }
  };

  // Sync state & trigger debounced fetch
  useEffect(() => {
    const timer = setTimeout(() => {
      fetchLiveLogs(searchQuery, level, service, host, since);
    }, 250);

    return () => clearTimeout(timer);
  }, [searchQuery, level, service, host, since]);

  // Compute live stats from currently active logs
  const metrics = useMemo(() => {
    const total = activeLogs.length;
    const errorCount = activeLogs.filter(log => log.level === 'error').length;
    const errorRate = total > 0 ? Math.round((errorCount / total) * 100) : 0;
    const activeServicesCount = new Set(activeLogs.map(log => log.service)).size;
    const activeHostsCount = new Set(activeLogs.map(log => log.host)).size;

    return {
      total,
      errorCount,
      errorRate,
      activeServicesCount,
      activeHostsCount
    };
  }, [activeLogs]);

  return (
    <div className="min-h-screen bg-[#060608] text-gray-100 p-6 md:p-8 font-sans selection:bg-indigo-500/30">
      <div className="max-w-7xl mx-auto space-y-8 animate-in fade-in duration-500">

        {/* Header Section */}
        <header className="flex flex-col lg:flex-row justify-between items-start lg:items-center gap-6 pb-6 border-b border-gray-800/40">
          <div>
            <div className="flex flex-wrap items-center gap-2.5">
              <h1 className="text-3xl font-extrabold text-white tracking-tight bg-clip-text bg-gradient-to-r from-white via-gray-100 to-indigo-200">
                Observability Hub
              </h1>
              <div className="flex items-center gap-1.5">
                <span className="flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                  <ShieldCheck className="h-3 w-3" />
                  Secure
                </span>
                
                {isLive ? (
                  <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-indigo-550/15 text-indigo-400 border border-indigo-500/20">
                    <span className="h-1.5 w-1.5 rounded-full bg-indigo-400 animate-pulse" />
                    Live Go API
                  </span>
                ) : (
                  <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-rose-500/10 text-rose-455 border border-rose-500/20">
                    <span className="h-1.5 w-1.5 rounded-full bg-rose-500" />
                    Offline
                  </span>
                )}
              </div>
            </div>
            <p className="text-sm text-gray-400 mt-1.5">Monitor, troubleshoot, and analyze system health telemetry in real-time.</p>
          </div>
          
          <SearchBar
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
            level={level}
            setLevel={setLevel}
            service={service}
            setService={setService}
            host={host}
            setHost={setHost}
            since={since}
            setSince={setSince}
            servicesList={SERVICES_LIST}
            hostsList={HOSTS_LIST}
          />
        </header>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {/* Total Logs Card */}
          <div className="bg-gray-950/40 border border-gray-855 p-5 rounded-2xl shadow-xl backdrop-blur-md space-y-2 flex flex-col justify-between hover:border-gray-800 transition-all duration-300">
            <span className="text-xs font-bold text-gray-455 uppercase tracking-wider flex items-center gap-1.5">
              <Terminal className="h-4 w-4 text-indigo-400" />
              Logs Returned
            </span>
            <div className="flex justify-between items-baseline mt-1">
              <span className="text-2xl font-extrabold text-white tracking-tight">{metrics.total}</span>
              <span className="text-[10px] font-semibold text-gray-455">Active Stream</span>
            </div>
            <div className="h-1.5 bg-gray-900 rounded-full overflow-hidden mt-2">
              <div className="h-full bg-indigo-500 rounded-full transition-all duration-500" style={{ width: '100%' }} />
            </div>
          </div>

          {/* Error Rate Card */}
          <div className="bg-gray-950/40 border border-gray-855 p-5 rounded-2xl shadow-xl backdrop-blur-md space-y-2 flex flex-col justify-between hover:border-gray-800 transition-all duration-300">
            <span className="text-xs font-bold text-gray-455 uppercase tracking-wider flex items-center gap-1.5">
              <Activity className="h-4 w-4 text-rose-455" />
              System Error Rate
            </span>
            <div className="flex justify-between items-baseline mt-1">
              <span className="text-2xl font-extrabold text-white tracking-tight">{metrics.errorRate}%</span>
              <span className="text-[10px] font-semibold text-rose-400">{metrics.errorCount} Failed</span>
            </div>
            <div className="h-1.5 bg-gray-900 rounded-full overflow-hidden mt-2">
              <div className={`h-full rounded-full transition-all duration-500 ${metrics.errorRate > 20 ? 'bg-rose-500 animate-pulse' : 'bg-emerald-500'}`} style={{ width: `${metrics.errorRate}%` }} />
            </div>
          </div>

          {/* Unique Services Card */}
          <div className="bg-gray-950/40 border border-gray-855 p-5 rounded-2xl shadow-xl backdrop-blur-md space-y-2 flex flex-col justify-between hover:border-gray-800 transition-all duration-300">
            <span className="text-xs font-bold text-gray-455 uppercase tracking-wider flex items-center gap-1.5">
              <Cpu className="h-4 w-4 text-purple-400" />
              Active Services
            </span>
            <div className="flex justify-between items-baseline mt-1">
              <span className="text-2xl font-extrabold text-white tracking-tight">{metrics.activeServicesCount}</span>
              <span className="text-[10px] font-semibold text-gray-450">Monitoring Online</span>
            </div>
            <div className="h-1.5 bg-gray-900 rounded-full overflow-hidden mt-2">
              <div className="h-full bg-purple-500 rounded-full transition-all duration-500" style={{ width: '100%' }} />
            </div>
          </div>

          {/* Unique Hosts Card */}
          <div className="bg-gray-950/40 border border-gray-855 p-5 rounded-2xl shadow-xl backdrop-blur-md space-y-2 flex flex-col justify-between hover:border-gray-800 transition-all duration-300">
            <span className="text-xs font-bold text-gray-455 uppercase tracking-wider flex items-center gap-1.5">
              <Server className="h-4 w-4 text-teal-400" />
              Monitored Hosts
            </span>
            <div className="flex justify-between items-baseline mt-1">
              <span className="text-2xl font-extrabold text-white tracking-tight">{metrics.activeHostsCount}</span>
              <span className="text-[10px] font-semibold text-gray-400">Instances Active</span>
            </div>
            <div className="h-1.5 bg-gray-900 rounded-full overflow-hidden mt-2">
              <div className="h-full bg-teal-500 rounded-full transition-all duration-500" style={{ width: '100%' }} />
            </div>
          </div>
        </div>

        {errorMsg && (
          <div className="bg-rose-500/10 border border-rose-500/25 p-4 rounded-xl flex items-start gap-3 text-rose-300">
            <AlertCircle className="h-5 w-5 shrink-0 mt-0.5" />
            <div className="space-y-1">
              <h4 className="text-sm font-bold">API Connection Error</h4>
              <p className="text-xs text-rose-300/80">{errorMsg}</p>
            </div>
          </div>
        )}

        <main className="space-y-6 relative">
          {/* Loader Overlay */}
          {isLoading && (
            <div className="absolute inset-0 bg-[#060608]/40 backdrop-blur-[1px] flex items-center justify-center z-10 rounded-2xl">
              <div className="p-3 bg-gray-955 border border-gray-800 rounded-full shadow-2xl flex items-center justify-center">
                <Loader2 className="h-6 w-6 text-indigo-400 animate-spin" />
              </div>
            </div>
          )}

          {/* Proactive Insights */}
          <section aria-label="System Insights">
            <IntentCard />
          </section>

          {/* Telemetry Data */}
          <section aria-label="Log Stream">
            <LogList logs={activeLogs} totalCount={activeLogs.length} />
          </section>
        </main>

      </div>
    </div>
  );
}
