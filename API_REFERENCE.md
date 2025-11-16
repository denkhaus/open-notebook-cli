# OpenNotebook API Reference

🔗 **Direct Links to Python API Source Code**

This document contains direct links to the Python API implementations for quick reference during development.

## 📁 API Structure

**Base Repository**: https://github.com/lfnovo/open-notebook/tree/main/api

## 🏗️ Service Architecture

### 🔄 HTTP Client
**File**: `api/client.py`
**URL**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/client.py
**Purpose**: Core HTTP client with request/response handling
**Key Info**: Uses `/_make_request()` method, handles both JSON and form data

### 📚 Individual Services

#### 🔗 Sources Service
**File**: `api/sources_service.py`
**URL**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/sources_service.py
**Purpose**: Source creation, listing, and management
**Key Methods**: `create_source()`, `get_sources()`, `get_source_status()`
**🚨 IMPORTANT**: Uses TWO different endpoints:
- **Web Interface**: `POST /api/sources` (multipart/form-data)
- **API Client**: `POST /api/sources/json` (application/json)

#### 💬 Chat Service
**File**: `api/chat_service.py`
**URL**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/chat_service.py
**Purpose**: Chat sessions, message handling, streaming
**Key Methods**: `create_session()`, `execute_chat()`, `stream_chat()`

#### 📓 Notebooks Service
**File**: `api/notebooks_service.py`
**URL**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/notebooks_service.py
**Purpose**: Notebook CRUD operations, source management
**Key Methods**: `create_notebook()`, `get_notebook()`, `add_source_to_notebook()`

#### 🤖 Models Service
**File**: `api/models_service.py`
**URL**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/models_service.py
**Purpose**: Model management, defaults, providers
**Key Methods**: `list_models()`, `get_default_models()`, `get_providers()`

#### 🔍 Search Service
**File**: `api/search_service.py`
**URL**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/search_service.py
**Purpose**: Search functionality, query handling
**Key Methods**: `search()`, `search_notebooks()`

#### ⚙️ Settings Service
**File**: `api/settings_service.py`
**URL**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/settings_service.py
**Purpose**: Application settings, configuration
**Key Methods**: `get_settings()`, `update_settings()`

#### 🔄 Command Service
**File**: `api/command_service.py`
**URL**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/command_service.py
**Purpose**: Background job management, async processing
**Key Methods**: `get_command_status()`, `cancel_command()`

## 🛣️ API Endpoints Patterns

### 🚨 Discovered Dual-Format Pattern
Some endpoints support **both** JSON and form data:

| Entity | JSON Endpoint | Form Data Endpoint | Usage |
|--------|---------------|-------------------|---------|
| Sources | `/api/sources/json` | `/api/sources` | ✅ Confirmed |
| Notebooks | `/api/notebooks/json` | `/api/notebooks` | ⚠️ Test needed |
| Models | `/api/models/json` | `/api/models` | ⚠️ Test needed |

### 📝 Standard Endpoints (JSON only)
- `/api/models` - Model listing and management
- `/api/models/defaults` - Default model settings
- `/api/models/providers` - Provider availability
- `/api/chat/sessions` - Chat session management
- `/api/chat/execute` - Chat execution
- `/api/notebooks` - Notebook management
- `/api/transformations` - Data transformations
- `/api/settings` - Application settings
- `/api/auth/status` - Authentication status

## 🔍 Development Guidelines

### 📋 Checklist for New Features

1. **Check Python API First**: Always check the corresponding service file before implementing
2. **Test Endpoint Formats**: Test both `/api/entity` and `/api/entity/json` endpoints
3. **Request Format Match**: Match the exact request format from Python implementation
4. **Field Validation**: Check required fields and validation rules from Python code

### 🚨 Known Issues

- **Sources Creation**: Must use `/api/sources/json` for JSON requests (NOT `/api/sources`)
- **Chat Sessions**: Requires `notebook_id` parameter
- **Authentication**: Bearer token required for some endpoints
- **Response Formats**: API may return arrays instead of wrapped objects

### 🎯 Quick Reference URLs

- **All API Files**: https://github.com/lfnovo/open-notebook/tree/main/api
- **Models/Types**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/models.py
- **Main Client**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/client.py
- **Router Definitions**: https://raw.githubusercontent.com/lfnovo/open-notebook/main/api/routers/

---

📝 **Last Updated**: 2025-11-14
🔄 **Update this document when new API patterns are discovered**