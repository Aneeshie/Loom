# Loom Web Dashboard Documentation

The Loom Web Dashboard is a telemetry interface built with **React**, **Next.js**, **Tailwind CSS v4**, and **Shadcn UI** components.

## Technical Design

### 1. State Coordination
All active state lives inside the top-level page component ([page.tsx](file:///Users/nara/loom/web/app/page.tsx)), ensuring sync between:
- The search text input.
- The advanced filter popover fields (Level, Service, Host, Since).
- The active filter badges.
- The metrics stats grid.
- The live log list stream.

### 2. Debouncing & API Optimization
To prevent excessive load on PostgreSQL and the Ollama LLM parser when a user types in the search bar, requests are debounced on the client side:
* **Debounce Delay**: `250ms` using React `useEffect` cleanups.
* **API Route Rewrites**: Configured in `next.config.ts` to proxy requests from the Next.js origin directly to the Go API gateway (`localhost:3000`) without triggering CORS issues.

### 3. Styling & Aesthetics
* **Theme**: Deep dark space theme using custom slate and indigo variables.
* **Glassmorphism**: Panels and grids utilize `backdrop-blur-xl bg-gray-950/40 border-gray-850/80` for a sleek modern look.
* **Micro-Animations**: Hover animations on stats cards, popover transitions, and blinking green live state indicators.

### 4. Interactive Components
* **Advanced Search Filters**: Implemented using a Radix-based Popover trigger button containing toggles and dropdown selectors.
* **Log Detail Inspector**: Slides out a sheet from the right showing metadata attributes, raw Unix timestamps, and syntax-highlighted context JSON with copy-to-clipboard actions.
