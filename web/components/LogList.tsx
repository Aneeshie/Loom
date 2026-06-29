"use client";

import { useState } from 'react';
import { LogEntry } from "@/lib/types";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Terminal, Copy, Check, Clock, Server, Cpu, FileJson, AlertCircle } from 'lucide-react';

interface LogListProps {
  logs: LogEntry[];
  totalCount: number;
}

const levelStyles = {
  info: 'text-emerald-400 bg-emerald-500/10 border-emerald-500/20 hover:bg-emerald-500/20',
  warn: 'text-amber-400 bg-amber-500/10 border-amber-500/20 hover:bg-amber-500/20',
  error: 'text-rose-400 bg-rose-500/10 border-rose-500/20 hover:bg-rose-500/20',
};

const levelBadgeStyles = {
  info: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20 hover:bg-emerald-500/20',
  warn: 'bg-amber-500/10 text-amber-400 border-amber-500/20 hover:bg-amber-500/20',
  error: 'bg-rose-500/10 text-rose-400 border-rose-500/20 hover:bg-rose-500/20',
};

export default function LogList({ logs, totalCount }: LogListProps) {
  const [selectedLog, setSelectedLog] = useState<LogEntry | null>(null);
  const [copiedType, setCopiedType] = useState<'message' | 'json' | null>(null);

  const handleCopy = (text: string, type: 'message' | 'json') => {
    navigator.clipboard.writeText(text).then(() => {
      setCopiedType(type);
      setTimeout(() => setCopiedType(null), 2000);
    });
  };

  const formatTimestamp = (isoString: string) => {
    try {
      const date = new Date(isoString);
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false }) + '.' + String(date.getMilliseconds()).padStart(3, '0');
    } catch {
      return isoString;
    }
  };

  const formatFullDate = (isoString: string) => {
    try {
      const date = new Date(isoString);
      return date.toLocaleString([], { dateStyle: 'medium', timeStyle: 'medium' });
    } catch {
      return isoString;
    }
  };

  return (
    <div className="bg-gray-950/40 border border-gray-800/80 rounded-2xl overflow-hidden shadow-2xl backdrop-blur-xl">
      {/* List Header */}
      <div className="px-6 py-4.5 border-b border-gray-800/60 bg-gray-900/30 flex justify-between items-center">
        <div className="flex items-center gap-3">
          <h2 className="text-base font-bold text-gray-100 tracking-wide flex items-center gap-2">
            <Terminal className="h-4.5 w-4.5 text-indigo-400" />
            Live Log Stream
          </h2>
          <span className="flex h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
        </div>
        <span className="text-xs text-gray-400 font-semibold bg-gray-900/80 border border-gray-800/60 px-2.5 py-1 rounded-lg">
          Showing {logs.length} of {totalCount} logs
        </span>
      </div>

      {/* Logs Table Content */}
      <div className="divide-y divide-gray-850/40 font-mono text-[13px] overflow-x-auto">
        {logs.length === 0 ? (
          <div className="px-6 py-16 flex flex-col items-center justify-center text-center gap-3 bg-gray-900/10">
            <div className="p-4 bg-gray-900/40 rounded-full border border-gray-800">
              <AlertCircle className="h-7 w-7 text-gray-500" />
            </div>
            <div className="space-y-1">
              <h3 className="text-sm font-semibold text-gray-200">No matching logs found</h3>
              <p className="text-xs text-gray-400 max-w-sm">
                Try refining your keyword query or clearing some of the advanced filter options.
              </p>
            </div>
          </div>
        ) : (
          logs.map((log) => (
            <div
              key={log.id}
              onClick={() => setSelectedLog(log)}
              className="px-6 py-3 hover:bg-indigo-500/[0.02] active:bg-indigo-500/[0.04] transition-all flex flex-col md:flex-row md:items-center gap-3 md:gap-6 group cursor-pointer border-l-2 border-l-transparent hover:border-l-indigo-500/50"
            >
              {/* Time Column */}
              <span className="text-gray-500 whitespace-nowrap group-hover:text-gray-400 transition-colors text-xs shrink-0 w-24">
                {formatTimestamp(log.timestamp)}
              </span>

              {/* Severity Level Column */}
              <span className={`px-2 py-0.5 rounded border uppercase text-[10px] font-bold tracking-wider w-16 text-center select-none shrink-0 ${levelStyles[log.level]}`}>
                {log.level}
              </span>

              {/* Service Tag Column */}
              <span className="text-purple-400 font-medium shrink-0 md:w-36 truncate flex items-center gap-1">
                <Cpu className="h-3 w-3 text-purple-500/70" />
                {log.service}
              </span>

              {/* Host Column (Desktop Only) */}
              <span className="hidden lg:flex items-center gap-1 text-teal-400/90 font-medium shrink-0 w-36 truncate">
                <Server className="h-3 w-3 text-teal-500/70" />
                {log.host}
              </span>

              {/* Message Column */}
              <span className="text-gray-300 flex-1 truncate group-hover:text-white transition-colors">
                {log.message}
              </span>
            </div>
          ))
        )}
      </div>

      {/* Log Detail Drawer (Sheet) */}
      <Sheet open={!!selectedLog} onOpenChange={(open) => !open && setSelectedLog(null)}>
        <SheetContent className="w-full sm:max-w-md bg-gray-950/98 border-l border-gray-800 text-gray-200 p-6 flex flex-col h-full shadow-2xl backdrop-blur-xl">
          {selectedLog && (
            <>
              <SheetHeader className="pb-4 border-b border-gray-800/80">
                <div className="flex items-center justify-between mb-2">
                  <span className={`px-2 py-0.5 rounded border uppercase text-[10px] font-bold tracking-wider ${levelBadgeStyles[selectedLog.level]}`}>
                    {selectedLog.level}
                  </span>
                  <span className="text-[11px] text-gray-500 font-mono">ID: {selectedLog.id}</span>
                </div>
                <SheetTitle className="text-base font-bold text-gray-100 flex items-center gap-2">
                  Log Details
                </SheetTitle>
                <SheetDescription className="text-xs text-gray-400">
                  Detailed inspection of telemetry event from service <strong>{selectedLog.service}</strong>.
                </SheetDescription>
              </SheetHeader>

              <div className="flex-1 overflow-y-auto py-6 space-y-6">
                {/* Event Message */}
                <div className="space-y-2">
                  <span className="text-[11px] font-bold text-gray-450 uppercase tracking-wider block">Log Message</span>
                  <div className="p-4 bg-gray-900/40 border border-gray-800/60 rounded-xl font-mono text-sm leading-relaxed text-gray-200 selection:bg-indigo-500/20 select-text">
                    {selectedLog.message}
                  </div>
                </div>

                {/* Context Metadata Cards */}
                <div className="grid grid-cols-2 gap-3">
                  <div className="p-3 bg-gray-900/20 border border-gray-800/50 rounded-xl space-y-1">
                    <span className="text-[10px] font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1">
                      <Cpu className="h-3 w-3 text-purple-400" />
                      Service
                    </span>
                    <span className="text-xs font-semibold text-gray-300 font-mono block truncate">{selectedLog.service}</span>
                  </div>
                  <div className="p-3 bg-gray-900/20 border border-gray-800/50 rounded-xl space-y-1">
                    <span className="text-[10px] font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1">
                      <Server className="h-3 w-3 text-teal-400" />
                      Host
                    </span>
                    <span className="text-xs font-semibold text-gray-300 font-mono block truncate">{selectedLog.host}</span>
                  </div>
                  <div className="col-span-2 p-3 bg-gray-900/20 border border-gray-800/50 rounded-xl space-y-1">
                    <span className="text-[10px] font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1">
                      <Clock className="h-3 w-3 text-indigo-400" />
                      Timestamp
                    </span>
                    <span className="text-xs font-semibold text-gray-300 font-mono block">
                      {formatFullDate(selectedLog.timestamp)}
                    </span>
                  </div>
                </div>

                {/* JSON Metadata Inspector */}
                {selectedLog.metadata && (
                  <div className="space-y-2 flex flex-col">
                    <span className="text-[11px] font-bold text-gray-450 uppercase tracking-wider block">Context Metadata</span>
                    <div className="relative flex-1 group">
                      <pre className="p-4 bg-gray-900/60 border border-gray-800/60 rounded-xl font-mono text-xs overflow-x-auto text-indigo-200/90 leading-normal select-text max-h-[300px]">
                        {JSON.stringify(selectedLog.metadata, null, 2)}
                      </pre>
                      
                      <button
                        onClick={() => handleCopy(JSON.stringify(selectedLog.metadata, null, 2), 'json')}
                        className="absolute top-3 right-3 p-1.5 bg-gray-950/80 border border-gray-800 text-gray-400 hover:text-white rounded-lg opacity-0 group-hover:opacity-100 focus:opacity-100 transition-all cursor-pointer"
                        title="Copy JSON Payload"
                      >
                        {copiedType === 'json' ? (
                          <Check className="h-3.5 w-3.5 text-emerald-400" />
                        ) : (
                          <Copy className="h-3.5 w-3.5" />
                        )}
                      </button>
                    </div>
                  </div>
                )}
              </div>

              {/* Action Buttons */}
              <div className="pt-4 border-t border-gray-800/80 flex gap-2.5">
                <Button
                  variant="outline"
                  size="sm"
                  className="flex-1 text-xs py-5 rounded-xl border-gray-800 hover:bg-gray-900 cursor-pointer"
                  onClick={() => handleCopy(selectedLog.message, 'message')}
                >
                  {copiedType === 'message' ? (
                    <>
                      <Check className="mr-1.5 h-3.5 w-3.5 text-emerald-400" />
                      Copied Msg
                    </>
                  ) : (
                    <>
                      <Copy className="mr-1.5 h-3.5 w-3.5" />
                      Copy Msg
                    </>
                  )}
                </Button>
                
                {selectedLog.metadata && (
                  <Button
                    size="sm"
                    className="flex-1 bg-indigo-600 hover:bg-indigo-500 text-white text-xs py-5 rounded-xl cursor-pointer"
                    onClick={() => handleCopy(JSON.stringify(selectedLog.metadata, null, 2), 'json')}
                  >
                    <FileJson className="mr-1.5 h-3.5 w-3.5" />
                    Copy JSON
                  </Button>
                )}
              </div>
            </>
          )}
        </SheetContent>
      </Sheet>
    </div>
  );
}
