"use client";

import { Zap, ArrowRight, AlertTriangle } from 'lucide-react';

export default function IntentCard() {
  return (
    <div className="relative group overflow-hidden rounded-2xl border border-indigo-500/30 bg-gradient-to-br from-indigo-950/40 via-slate-900/50 to-blue-950/30 p-6 shadow-2xl backdrop-blur-xl transition-all duration-300 hover:shadow-indigo-500/5">
      {/* Dynamic light streak animation */}
      <div className="absolute -inset-y-12 -inset-x-6 w-[200%] bg-gradient-to-r from-transparent via-white/5 to-transparent skew-x-12 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000 ease-in-out" />

      {/* Background glow orb */}
      <div className="absolute top-0 right-0 -mt-8 -mr-8 h-32 w-32 rounded-full bg-indigo-500/10 blur-2xl group-hover:bg-indigo-500/15 transition-all duration-500" />

      <div className="relative flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div className="flex items-start space-x-4">
          <div className="p-3 bg-indigo-500/20 rounded-xl shadow-inner border border-indigo-400/25 shrink-0 flex items-center justify-center animate-pulse">
            <Zap className="h-6 w-6 text-indigo-400 fill-indigo-400/10" />
          </div>
          <div className="space-y-1.5">
            <div className="flex flex-wrap items-center gap-2">
              <h3 className="text-base font-bold text-white tracking-wide">
                Proactive Insight: High Error Rate Detected
              </h3>
              <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-rose-500/10 text-rose-400 border border-rose-500/20">
                <AlertTriangle className="h-2.5 w-2.5" />
                Alert
              </span>
            </div>

            <p className="text-sm text-gray-300 leading-relaxed max-w-3xl">
              The service <code className="bg-gray-950/80 border border-gray-800/80 px-1.5 py-0.5 rounded text-indigo-300 font-mono text-xs">payment-gateway</code> is experiencing a <span className="text-rose-400 font-semibold">15.2% spike</span> in 5xx error statuses over the past 10 minutes.
              This anomaly correlates directly with deployment version <code className="bg-gray-950/80 border border-gray-800/80 px-1.5 py-0.5 rounded text-indigo-300 font-mono text-xs">v1.4.2</code> on host <code className="bg-gray-950/80 border border-gray-800/80 px-1.5 py-0.5 rounded text-indigo-300 font-mono text-xs">prod-ap-south-1a</code>.
            </p>
          </div>
        </div>

        <button className="shrink-0 flex items-center justify-center px-4.5 py-2.5 text-xs font-bold text-white bg-indigo-600 hover:bg-indigo-500 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 focus:ring-offset-gray-955 transition-all duration-200 shadow-lg shadow-indigo-600/20 active:scale-98 cursor-pointer">
          Investigate Incident
          <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-0.5" />
        </button>
      </div>
    </div>
  );
}
