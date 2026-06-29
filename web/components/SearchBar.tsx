"use client";

import { useEffect, useRef, useState } from 'react';
import { Search, SlidersHorizontal, X, RotateCcw, Clock, Server, Cpu, Activity } from 'lucide-react';
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

interface SearchBarProps {
  searchQuery: string;
  setSearchQuery: (val: string) => void;
  level: string;
  setLevel: (val: string) => void;
  service: string;
  setService: (val: string) => void;
  host: string;
  setHost: (val: string) => void;
  since: string;
  setSince: (val: string) => void;
  servicesList: string[];
  hostsList: string[];
}

export default function SearchBar({
  searchQuery,
  setSearchQuery,
  level,
  setLevel,
  service,
  setService,
  host,
  setHost,
  since,
  setSince,
  servicesList,
  hostsList,
}: SearchBarProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [isOpen, setIsOpen] = useState(false);

  // Focus search bar on Command+K or "/"
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        inputRef.current?.focus();
      } else if (e.key === '/' && document.activeElement !== inputRef.current) {
        // Prevent default only if not typing in another input
        const activeEl = document.activeElement;
        const isInput = activeEl?.tagName === 'INPUT' || activeEl?.tagName === 'TEXTAREA' || activeEl?.hasAttribute('contenteditable');
        if (!isInput) {
          e.preventDefault();
          inputRef.current?.focus();
        }
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  const hasActiveFilters = level !== 'all' || service !== '' || host !== '' || since !== 'all';

  const resetFilters = () => {
    setLevel('all');
    setService('');
    setHost('');
    setSince('all');
  };

  const getSinceLabel = (val: string) => {
    switch (val) {
      case '5m': return 'Last 5 mins';
      case '15m': return 'Last 15 mins';
      case '1h': return 'Last 1 hour';
      case '24h': return 'Last 24 hours';
      case '7d': return 'Last 7 days';
      default: return 'All time';
    }
  };

  return (
    <div className="w-full max-w-3xl flex flex-col gap-3">
      {/* Search Input Box */}
      <div className="relative flex items-center w-full group">
        {/* Glow background effect */}
        <div className="absolute -inset-0.5 bg-gradient-to-r from-indigo-500/20 via-purple-500/10 to-blue-500/20 rounded-xl blur opacity-30 group-focus-within:opacity-100 transition duration-500" />
        
        <div className="relative flex items-center w-full bg-gray-900/90 border border-gray-800 rounded-xl shadow-2xl backdrop-blur-md overflow-hidden focus-within:border-indigo-500/50 focus-within:ring-1 focus-within:ring-indigo-500/50 transition-all duration-300">
          <div className="pl-3.5 flex items-center pointer-events-none text-gray-500">
            <Search className="h-5 w-5 transition-colors group-focus-within:text-indigo-400" />
          </div>
          
          <input
            ref={inputRef}
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="flex-1 bg-transparent px-3 py-3.5 text-sm text-gray-100 placeholder-gray-500 focus:outline-none"
            placeholder="Search logs by keyword, exception, message, or use filters..."
          />
          
          {searchQuery && (
            <button
              onClick={() => setSearchQuery('')}
              className="p-1 mr-1 text-gray-500 hover:text-gray-300 rounded-lg hover:bg-gray-800/50 transition-colors"
            >
              <X className="h-4 w-4" />
            </button>
          )}

          {/* Kbd shortcut indicator */}
          <div className="hidden sm:flex items-center text-[10px] text-gray-500 bg-gray-950/80 px-2 py-1 rounded-md border border-gray-800/80 font-mono gap-0.5 select-none pointer-events-none mr-2">
            <span>⌘</span>
            <span>K</span>
          </div>

          <div className="h-6 w-px bg-gray-800 self-center mr-1" />

          {/* Advanced Filter Popover */}
          <Popover open={isOpen} onOpenChange={setIsOpen}>
            <PopoverTrigger asChild>
              <Button
                variant="ghost"
                size="sm"
                className={`h-9 px-3 mr-1.5 flex items-center gap-1.5 rounded-lg transition-all ${
                  hasActiveFilters 
                    ? 'text-indigo-400 hover:text-indigo-300 bg-indigo-500/10 hover:bg-indigo-500/20' 
                    : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800/80'
                }`}
              >
                <SlidersHorizontal className="h-4 w-4" />
                <span className="hidden sm:inline text-xs font-semibold">Filters</span>
                {hasActiveFilters && (
                  <span className="flex h-2 w-2 rounded-full bg-indigo-400 animate-pulse" />
                )}
              </Button>
            </PopoverTrigger>
            
            <PopoverContent 
              align="end" 
              className="w-80 bg-gray-950/95 border border-gray-850 shadow-2xl rounded-xl p-4 backdrop-blur-xl z-50 text-gray-200 gap-4 flex flex-col"
            >
              <div className="flex justify-between items-center pb-2 border-b border-gray-800/60">
                <span className="font-semibold text-sm tracking-wide text-gray-100 flex items-center gap-1.5">
                  <SlidersHorizontal className="h-4 w-4 text-indigo-400" />
                  Advanced Filters
                </span>
                {hasActiveFilters && (
                  <button 
                    onClick={resetFilters}
                    className="text-[11px] text-gray-400 hover:text-rose-450 transition-colors flex items-center gap-1"
                  >
                    <RotateCcw className="h-3 w-3" />
                    Reset
                  </button>
                )}
              </div>

              {/* LEVEL FILTER */}
              <div className="space-y-2">
                <label className="text-[11px] font-bold text-gray-450 uppercase tracking-wider flex items-center gap-1">
                  <Activity className="h-3.5 w-3.5 text-emerald-400" />
                  Severity Level
                </label>
                <div className="grid grid-cols-4 gap-1">
                  {['all', 'info', 'warn', 'error'].map((l) => (
                    <button
                      key={l}
                      onClick={() => setLevel(l)}
                      className={`py-1.5 text-xs font-semibold rounded-md border transition-all capitalize cursor-pointer ${
                        level === l
                          ? l === 'error'
                            ? 'bg-rose-500/20 text-rose-300 border-rose-500/50'
                            : l === 'warn'
                            ? 'bg-amber-500/20 text-amber-300 border-amber-500/50'
                            : l === 'info'
                            ? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/50'
                            : 'bg-indigo-500/20 text-indigo-300 border-indigo-500/50'
                          : 'bg-gray-900/60 text-gray-405 border-gray-800/80 hover:bg-gray-800/40 hover:text-gray-300'
                      }`}
                    >
                      {l}
                    </button>
                  ))}
                </div>
              </div>

              {/* SERVICE FILTER */}
              <div className="space-y-2">
                <label className="text-[11px] font-bold text-gray-450 uppercase tracking-wider flex items-center gap-1">
                  <Cpu className="h-3.5 w-3.5 text-purple-400" />
                  Service
                </label>
                <select
                  value={service}
                  onChange={(e) => setService(e.target.value)}
                  className="w-full bg-gray-900/80 border border-gray-800/85 rounded-lg px-3 py-2 text-xs text-gray-200 focus:outline-none focus:border-indigo-500/60 cursor-pointer"
                >
                  <option value="">All Services</option>
                  {servicesList.map((srv) => (
                    <option key={srv} value={srv}>{srv}</option>
                  ))}
                </select>
              </div>

              {/* HOST FILTER */}
              <div className="space-y-2">
                <label className="text-[11px] font-bold text-gray-450 uppercase tracking-wider flex items-center gap-1">
                  <Server className="h-3.5 w-3.5 text-teal-400" />
                  Host
                </label>
                <select
                  value={host}
                  onChange={(e) => setHost(e.target.value)}
                  className="w-full bg-gray-900/80 border border-gray-800/85 rounded-lg px-3 py-2 text-xs text-gray-200 focus:outline-none focus:border-indigo-500/60 cursor-pointer"
                >
                  <option value="">All Hosts</option>
                  {hostsList.map((h) => (
                    <option key={h} value={h}>{h}</option>
                  ))}
                </select>
              </div>

              {/* SINCE FILTER */}
              <div className="space-y-2">
                <label className="text-[11px] font-bold text-gray-450 uppercase tracking-wider flex items-center gap-1">
                  <Clock className="h-3.5 w-3.5 text-indigo-400" />
                  Time Window
                </label>
                <div className="grid grid-cols-3 gap-1">
                  {['all', '5m', '15m', '1h', '24h', '7d'].map((s) => (
                    <button
                      key={s}
                      onClick={() => setSince(s)}
                      className={`py-1.5 text-xs font-semibold rounded-md border transition-all cursor-pointer ${
                        since === s
                          ? 'bg-indigo-500/20 text-indigo-300 border-indigo-500/50'
                          : 'bg-gray-900/60 text-gray-405 border-gray-800/80 hover:bg-gray-800/40 hover:text-gray-300'
                      }`}
                    >
                      {s === 'all' ? 'All time' : s}
                    </button>
                  ))}
                </div>
              </div>

              <div className="pt-2 border-t border-gray-900 flex justify-end">
                <Button 
                  size="sm"
                  className="bg-indigo-600 hover:bg-indigo-500 text-white font-semibold rounded-lg text-xs px-4 cursor-pointer"
                  onClick={() => setIsOpen(false)}
                >
                  Done
                </Button>
              </div>
            </PopoverContent>
          </Popover>
        </div>
      </div>

      {/* Active Filter Badges */}
      {hasActiveFilters && (
        <div className="flex flex-wrap items-center gap-2 px-1 text-xs text-gray-400 animate-in fade-in slide-in-from-top-1 duration-300">
          <span className="text-[11px] text-gray-500">Active filters:</span>
          
          {level !== 'all' && (
            <Badge 
              variant="outline" 
              className={`capitalize px-2 py-0.5 rounded-lg flex items-center gap-1 text-[11px] ${
                level === 'error' ? 'text-rose-450 border-rose-500/30 bg-rose-500/5' :
                level === 'warn' ? 'text-amber-455 border-amber-500/30 bg-amber-500/5' :
                'text-emerald-450 border-emerald-500/30 bg-emerald-500/5'
              }`}
            >
              Level: {level}
              <button onClick={() => setLevel('all')} className="hover:text-white transition-colors cursor-pointer">
                <X className="h-3 w-3" />
              </button>
            </Badge>
          )}

          {service !== '' && (
            <Badge 
              variant="outline" 
              className="text-purple-400 border-purple-500/30 bg-purple-500/5 px-2 py-0.5 rounded-lg flex items-center gap-1 text-[11px]"
            >
              Service: {service}
              <button onClick={() => setService('')} className="hover:text-white transition-colors cursor-pointer">
                <X className="h-3 w-3" />
              </button>
            </Badge>
          )}

          {host !== '' && (
            <Badge 
              variant="outline" 
              className="text-teal-400 border-teal-500/30 bg-teal-500/5 px-2 py-0.5 rounded-lg flex items-center gap-1 text-[11px]"
            >
              Host: {host}
              <button onClick={() => setHost('')} className="hover:text-white transition-colors cursor-pointer">
                <X className="h-3 w-3" />
              </button>
            </Badge>
          )}

          {since !== 'all' && (
            <Badge 
              variant="outline" 
              className="text-indigo-400 border-indigo-500/30 bg-indigo-500/5 px-2 py-0.5 rounded-lg flex items-center gap-1 text-[11px]"
            >
              Since: {getSinceLabel(since)}
              <button onClick={() => setSince('all')} className="hover:text-white transition-colors cursor-pointer">
                <X className="h-3 w-3" />
              </button>
            </Badge>
          )}

          <button 
            onClick={resetFilters}
            className="text-[10px] text-gray-500 hover:text-rose-400 transition-colors ml-1 underline underline-offset-2 flex items-center gap-0.5 cursor-pointer"
          >
            Clear all
          </button>
        </div>
      )}
    </div>
  );
}
