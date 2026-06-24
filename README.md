# Schedulord

Schedulord is a simple, configurable job scheduler written in Go, designed to run as a containerized service. It provides cron-based triggering and maps schedules to registered handlers, enabling deterministic execution of scheduled jobs within a lightweight, minimal runtime.

---

## Overview

The system is designed around predictable execution of scheduled tasks in constrained environments. It focuses on determinism, low overhead, and deployment simplicity.

Key characteristics:

- Cron-based scheduling model for recurring execution.
- Handler-based architecture for mapping jobs to executable logic.
- Time-zone aware scheduling using IANA time zones.
- GitHub API integration via Personal Access Tokens (optional execution target).
- Graceful shutdown handling for controlled termination in container environments.
- Static binary distribution optimized for minimal container images.

---

## Architecture

Schedulord is structured as a modular Go application with clear separation of concerns:

- **Entry layer (`cmd/`)**  
  Application bootstrap and runtime initialization.

- **Configuration layer (`internal/config/`)**  
  Loads and validates runtime configuration.

- **Scheduling engine (`internal/scheduler/`)**  
  Core cron evaluation and job dispatch logic.

- **Integration layer**  
  External system adapters (e.g., GitHub API execution triggers).

---

## Design Goals

- Deterministic execution of scheduled tasks.
- Minimal runtime overhead suitable for containers.
- Clear separation between scheduling logic and execution handlers.
- Config-driven behavior without hardcoded scheduling rules.
- Portable deployment via static Go binary.

---

## Status

This is a private, source-available project intended for demonstration of engineering capability and system design.