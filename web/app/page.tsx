"use client";

import { useState, useMemo } from 'react';
import SearchBar from '@/components/SearchBar';
import IntentCard from '@/components/IntentCard';
import LogList from '@/components/LogList';
import { getMockLogs } from '@/lib/mockLogs';
import { Terminal, Activity, Cpu, Server, ShieldCheck } from 'lucide-react';

export default function Dashboard() {
  // Main filter states
  const [searchQuery, setSearchQuery] = useState('');
  const [level, setLevel] = useState('all');
  const [service, setService] = useState('');
  const [host, setHost] = useState('');
  const [since, setSince] = useState('all');

  // Load static logs once
  const logs = useMemo(() => getMockLogs(), []);

  // Compute unique lists for filter selectors
  const servicesList = useMemo(() => {
    return Array.from(new Set(logs.map(log => log.service))).sort();
  }, [logs]);

  const hostsList = useMemo(() => {
    return Array.from(new Set(logs.map(log => log.host))).sort();
  }, [logs]);

  // Compute filtered logs
  const filteredLogs = useMemo(() => {
    return logs.filter(log => {
      // 1. Level filter
      if (level !== 'all' && log.level !== level.toLowerCase()) {
        return false;
      }
      
      // 2. Service filter
      if (service && log.service !== service) {
        return false;
      }
      
      // 3. Host filter
      if (host && log.host !== host) {
        return false;
      }
      
      // 4. Since (time window) filter
      if (since !== 'all') {
        const logTime = new Date(log.timestamp).getTime();
        const nowTime = Date.now();
        let thresholdMs = 0;
        
        if (since === '5m') thresholdMs = 5 * 60 * 1000;
        else if (since === '15m') thresholdMs = 15 * 60 * 1000;
        else if (since === '1h') thresholdMs = 60 * 60 * 1000;
        else if (since === '24h') thresholdMs = 24 * 60 * 60 * 1000;
        else if (since === '7d') thresholdMs = 7 * 24 * 60 * 60 * 1000;
        
        if (nowTime - logTime > thresholdMs) {
          return false;
        }
      }
      
      // 5. Search query keyword filter
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        const matchesMessage = log.message.toLowerCase().includes(q);
        const matchesService = log.service.toLowerCase().includes(q);
        const matchesHost = log.host.toLowerCase().includes(q);
        const matchesLevel = log.level.toLowerCase().includes(q);
        
        return matchesMessage || matchesService || matchesHost || matchesLevel;
      }
      
      return true;
    });
  }, [logs, searchQuery, level, service, host, since]);

  // Compute telemetry metrics for stats panel
  const metrics = useMemo(() => {
    const total = filteredLogs.length;
    const errorCount = filteredLogs.filter(log => log.level === 'error').length;
    const errorRate = total > 0 ? Math.round((errorCount / total) * 100) : 0;
    const activeServicesCount = new Set(filteredLogs.map(log => log.service)).size;
    const activeHostsCount = new Set(filteredLogs.map(log => log.host)).size;

    return {
      total,
      errorCount,
      errorRate,
      activeServicesCount,
      activeHostsCount
    };
  }, [filteredLogs]);

  return (
    <div className="min-h-screen bg-[#060608] text-gray-100 p-6 md:p-8 font-sans selection:bg-indigo-500/30">
      <div className="max-w-7xl mx-auto space-y-8 animate-in fade-in duration-500">

        {/* Header Section */}
        <header className="flex flex-col lg:flex-row justify-between items-start lg:items-center gap-6 pb-6 border-b border-gray-800/40">
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-3xl font-extrabold text-white tracking-tight bg-clip-text bg-gradient-to-r from-white via-gray-100 to-indigo-200">
                Observability Hub
              </h1>
              <span className="flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                <ShieldCheck className="h-3 w-3" />
                Secure
              </span>
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
            servicesList={servicesList}
            hostsList={hostsList}
          />
        </header>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {/* Total Logs Card */}
          <div className="bg-gray-950/40 border border-gray-850 p-5 rounded-2xl shadow-xl backdrop-blur-md space-y-2 flex flex-col justify-between hover:border-gray-800 transition-all duration-300">
            <span className="text-xs font-bold text-gray-450 uppercase tracking-wider flex items-center gap-1.5">
              <Terminal className="h-4 w-4 text-indigo-400" />
              Logs Filtered
            </span>
            <div className="flex justify-between items-baseline mt-1">
              <span className="text-2xl font-extrabold text-white tracking-tight">{metrics.total}</span>
              <span className="text-[10px] font-semibold text-gray-450">Total: {logs.length}</span>
            </div>
            <div className="h-1.5 bg-gray-900 rounded-full overflow-hidden mt-2">
              <div className="h-full bg-indigo-500 rounded-full transition-all duration-500" style={{ width: `${(metrics.total / Math.max(logs.length, 1)) * 100}%` }} />
            </div>
          </div>

          {/* Error Rate Card */}
          <div className="bg-gray-950/40 border border-gray-850 p-5 rounded-2xl shadow-xl backdrop-blur-md space-y-2 flex flex-col justify-between hover:border-gray-800 transition-all duration-300">
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
          <div className="bg-gray-950/40 border border-gray-850 p-5 rounded-2xl shadow-xl backdrop-blur-md space-y-2 flex flex-col justify-between hover:border-gray-800 transition-all duration-300">
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
          <div className="bg-gray-950/40 border border-gray-850 p-5 rounded-2xl shadow-xl backdrop-blur-md space-y-2 flex flex-col justify-between hover:border-gray-800 transition-all duration-300">
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

        <main className="space-y-6">
          {/* Proactive Insights */}
          <section aria-label="System Insights">
            <IntentCard />
          </section>

          {/* Telemetry Data */}
          <section aria-label="Log Stream">
            <LogList logs={filteredLogs} totalCount={logs.length} />
          </section>
        </main>

      </div>
    </div>
  );
}
