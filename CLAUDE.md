# Claude AI - Project Workflow Guide

## 🎯 Project Overview

**OpenNotebook CLI** - Go CLI application for controlling OpenNotebook API (localhost:5055)

- Repo: `github.com/denkhaus/open-notebook-cli`
- Tech Stack: Go, urfav/cli/v2, samber/do/v2, zap-logger
- Architecture: HTTP Client with DI for OpenNotebook API
- All information about the API, like models, routes MUST be explored from [sourcecode](https://raw.githubusercontent.com/lfnovo/open-notebook/main/api)
- A Reference for each entity to the original python API can be found here @API_REFERENCE.md

---

## 🧠⚡ **INTEGRATED WORKFLOW SYSTEM**

### **Dual-Layer Architecture**

- **KNOT Tool**: **ONLY Source of Truth** for task management (Drop-in Replacement for Todo Tool)
- **Brain Tools**: Memory layer for knowledge management, user decisions, and project context
