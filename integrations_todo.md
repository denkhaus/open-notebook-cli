# Integration Tests Coverage Analysis & TODO

## Overview

This document outlines the current state of integration test coverage for the OpenNotebook CLI application and identifies entities/features that still need integration testing.

## Current Integration Test Coverage

### ✅ **Well Covered Entities:**

#### **Notes Operations** (`test/integration/notes_search_test.go`)
- ✅ CRUD operations (create, read, update, delete)
- ✅ Human note creation with notebook association
- ✅ Note content validation
- ✅ Note type validation (human/AI)

#### **Search Operations** (`test/integration/notes_search_test.go`)
- ✅ Text search functionality
- ✅ Vector search with relevance scoring
- ✅ Ask simple questions with AI models
- ✅ Search filtering (sources vs notes)
- ✅ Empty query handling
- ✅ Search result validation

#### **Notebooks** (`test/integration/simple_test.go`, `test/integration/client_test.go`)
- ✅ Basic CRUD operations
- ✅ Notebook listing and pagination
- ✅ Notebook creation and deletion
- ✅ Notebook metadata validation

#### **Authentication** (`test/integration/simple_test.go`, `test/integration/client_test.go`)
- ✅ Auth flow and password validation
- ✅ Auth status endpoint testing
- ✅ Authenticated vs unauthenticated access

#### **Models Management** (`test/integration/client_test.go`)
- ✅ Model listing from API
- ✅ Model type validation (language, embedding, TTS, STT)
- ✅ Model provider information
- ✅ Type-safe enum validation

#### **Settings** (`test/integration/simple_test.go`, `test/integration/client_test.go`)
- ✅ Settings endpoint access
- ✅ Settings model validation
- ✅ Enum types validation (YesNoDecision, ContentProcessingEngine, etc.)

#### **HTTP/Network Infrastructure** (Multiple test files)
- ✅ Error handling and retry logic
- ✅ Timeout management
- ✅ Network connectivity testing
- ✅ Response structure validation
- ✅ Context cancellation

#### **Streaming Operations** (`test/integration/streaming_simple_test.go`)
- ✅ Server-Sent Events (SSE) parsing
- ✅ Streaming response formats
- ✅ Context cancellation during streaming
- ✅ Performance characteristics

#### **Basic Sources** (`test/integration/client_test.go`)
- ✅ Source listing
- ✅ Text source creation
- ✅ Basic source validation

---

## 🚨 **Missing Integration Test Coverage**

### **1. Chat Operations** - 🔥 **CRITICAL PRIORITY**

**Available Models:**
- `ChatSession`, `ChatMessage`, `ChatExecuteRequest`, `ChatExecuteResponse`
- `ChatContextRequest`, `ChatCreateRequest`

**Available Commands:**
- `chat sessions list/create/delete`
- `chat execute`
- `chat history`

**Missing Test Coverage:**
- [ ] Chat session creation and management
- [ ] Chat message execution (streaming and non-streaming)
- [ ] Chat history retrieval
- [ ] Context-aware chat with notebook/source integration
- [ ] Multi-turn conversation handling
- [ ] Model switching during conversations
- [ ] Chat session persistence

**Business Impact:** HIGH - Chat is a core user-facing feature for AI interaction

---

### **2. Podcast Generation** - 🔥 **CRITICAL PRIORITY**

**Available Models:**
- `PodcastGenerationRequest`, `PodcastEpisode`, `PodcastEpisodeResponse`
- `PodcastJobStatus`, `PodcastEpisodesListResponse`

**Available Commands:**
- `podcast generate`
- `podcast episodes list/show/download/delete`

**Missing Test Coverage:**
- [ ] Podcast generation from sources
- [ ] Podcast generation from notebooks
- [ ] Podcast generation with custom queries
- [ ] Episode listing and pagination
- [ ] Episode details retrieval
- [ ] Audio file download functionality
- [ ] Episode deletion
- [ ] Job status tracking during generation
- [ ] Voice and language parameter validation
- [ ] Style parameter handling

**Business Impact:** HIGH - Unique content generation feature

---

### **3. Jobs & Background Operations** - 🔥 **HIGH PRIORITY**

**Available Models:**
- `JobStatus`, `JobsListResponse`

**Available Commands:**
- `jobs list`
- `jobs status <job-id>`
- `jobs cancel <job-id>`

**Missing Test Coverage:**
- [ ] Job listing with status filtering
- [ ] Job status monitoring and updates
- [ ] Job cancellation
- [ ] Job progress tracking
- [ ] Long-running operation handling
- [ ] Job pagination and limits
- [ ] Job error handling and failed status

**Business Impact:** HIGH - Critical for monitoring async operations

---

### **4. Transformations** - 🔥 **HIGH PRIORITY**

**Available Models:**
- `TransformationCreate`, `TransformationUpdate`, `Transformation`
- `TransformationExecuteRequest`, `TransformationExecuteResponse`
- `DefaultPromptResponse`, `DefaultPromptUpdate`

**Available Commands:**
- `transformations list/create/show/update/delete`
- `transformations execute`

**Missing Test Coverage:**
- [ ] Transformation CRUD operations
- [ ] Transformation execution on text
- [ ] Model selection for transformations
- [ ] Streaming transformation execution
- [ ] Default transformation management
- [ ] Prompt template validation
- [ ] Apply-default functionality

**Business Impact:** MEDIUM-HIGH - Important for content processing workflows

---

### **5. Advanced Source Operations** - 🔶 **MEDIUM PRIORITY**

**Missing Test Coverage:**
- [ ] File upload sources (currently only text sources tested)
- [ ] Link-based sources with URL processing
- [ ] Async source processing
- [ ] Source transformation application
- [ ] Source embedding process
- [ ] Source status tracking (pending, running, completed, failed)
- [ ] Source content processing engines (auto, docling, simple)
- [ ] Source deletion with files
- [ ] Batch source operations

**Business Impact:** MEDIUM - Important for content ingestion workflows

---

### **6. Embeddings & Rebuild Operations** - 🔶 **MEDIUM PRIORITY**

**Available Models:**
- `EmbedRequest`, `EmbedResponse`
- `RebuildRequest`, `RebuildResponse`, `RebuildStatusResponse`
- `RebuildProgress`, `RebuildStats`

**Missing Test Coverage:**
- [ ] Manual embedding of sources
- [ ] Manual embedding of notes
- [ ] Async embedding processing
- [ ] Knowledge base rebuilding (existing mode)
- [ ] Knowledge base rebuilding (all mode)
- [ ] Rebuild progress monitoring
- [ ] Rebuild statistics validation
- [ ] Rebuild error handling
- [ ] Include/exclude options (sources, notes, insights)

**Business Impact:** MEDIUM - Important for search functionality maintenance

---

### **7. Insights Generation** - 🔶 **MEDIUM PRIORITY**

**Available Models:**
- `SourceInsightResponse`, `CreateSourceInsightRequest`
- `SaveAsNoteRequest`

**Missing Test Coverage:**
- [ ] Source insight generation (summary type)
- [ ] Source insight generation (analysis type)
- [ ] Source insight generation (extraction type)
- [ ] Source insight generation (question type)
- [ ] Source insight generation (reflection type)
- [ ] Insight saving as notes
- [ ] Model selection for insights
- [ ] Transformation-based insight generation

**Business Impact:** MEDIUM - Enhances content analysis capabilities

---

### **8. Context Management** - 🔶 **LOWER PRIORITY**

**Available Models:**
- `ContextRequest`, `ContextResponse`, `ContextConfig`
- `ContextLevel` enum (low, medium, high, critical)

**Missing Test Coverage:**
- [ ] Context configuration for notebooks
- [ ] Context level assignment (low, medium, high, critical)
- [ ] Context token counting
- [ ] Context retrieval and validation
- [ ] Source-specific context levels
- [ ] Note-specific context levels

**Business Impact:** LOW-MEDIUM - Advanced feature for power users

---

## **Implementation Recommendations**

### **Phase 1: Critical Features** (Immediate Priority)
1. **Chat Operations Integration Tests**
   - File: `test/integration/chat_operations_test.go`
   - Focus: Session management, message execution, context handling

2. **Podcast Generation Integration Tests**
   - File: `test/integration/podcast_generation_test.go`
   - Focus: Generation workflow, episode management, job tracking

3. **Jobs Management Integration Tests**
   - File: `test/integration/jobs_management_test.go`
   - Focus: Job lifecycle, status monitoring, cancellation

### **Phase 2: Content Processing** (Secondary Priority)
4. **Transformations Integration Tests**
   - File: `test/integration/transformations_test.go`
   - Focus: CRUD operations, execution, prompt management

5. **Advanced Source Operations Integration Tests**
   - File: `test/integration/advanced_sources_test.go`
   - Focus: File uploads, URL processing, async operations

### **Phase 3: Advanced Features** (Future Priority)
6. **Embeddings & Rebuild Integration Tests**
   - File: `test/integration/embeddings_rebuild_test.go`
   - Focus: Embedding process, rebuild operations, progress tracking

7. **Insights Generation Integration Tests**
   - File: `test/integration/insights_generation_test.go`
   - Focus: Different insight types, saving workflows

8. **Context Management Integration Tests**
   - File: `test/integration/context_management_test.go`
   - Focus: Context configuration, level management

---

## **Test Implementation Guidelines**

### **Test Structure Standards**
```go
// Follow existing pattern from notes_search_test.go
func TestEntityOperations(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping entity integration tests")
    }
    
    if !isAPIAvailable("http://localhost:5055") {
        t.Skip("API not available on localhost:5055")
    }
    
    t.Run("Basic operations", func(t *testing.T) {
        // Test basic CRUD
    })
    
    t.Run("Advanced operations", func(t *testing.T) {
        // Test complex workflows
    })
    
    t.Run("Error scenarios", func(t *testing.T) {
        // Test error handling
    })
}
```

### **Authentication Handling**
- Use consistent auth setup from existing tests
- Handle auth failures gracefully
- Test both authenticated and unauthenticated scenarios where applicable

### **Cleanup Requirements**
- Always clean up test data (sessions, episodes, transformations, jobs)
- Use defer statements for cleanup
- Handle cleanup failures gracefully

### **Async Operation Testing**
- Implement polling for job status
- Set reasonable timeouts
- Test cancellation scenarios
- Validate progress updates

### **API Availability**
- Use `isAPIAvailable()` helper function
- Skip tests gracefully when API is unavailable
- Provide clear instructions for running with live API

---

## **Expected Benefits**

### **Quality Improvements**
- **Regression Detection**: Catch breaking changes in critical user workflows
- **API Contract Validation**: Ensure model structures match API responses
- **Error Handling**: Validate error scenarios and user experience
- **Performance Monitoring**: Track response times for key operations

### **Development Velocity**
- **Confidence in Changes**: Safe refactoring and feature additions
- **Documentation**: Tests serve as living documentation
- **Onboarding**: New developers understand workflows through tests
- **CI/CD**: Automated validation of deployments

### **User Experience**
- **Feature Completeness**: Ensure all advertised features work end-to-end
- **Cross-feature Integration**: Validate workflows that span multiple entities
- **Real-world Scenarios**: Test complex user journeys

---

## **Estimated Implementation Effort**

| Priority | Test Suite | Estimated Effort | Dependencies |
|----------|------------|------------------|--------------|
| 🔥 Critical | Chat Operations | 2-3 days | API running, auth setup |
| 🔥 Critical | Podcast Generation | 2-3 days | Audio processing, job system |
| 🔥 High | Jobs Management | 1-2 days | Background job system |
| 🔥 High | Transformations | 2 days | Model system, prompt templates |
| 🔶 Medium | Advanced Sources | 2 days | File system, URL processing |
| 🔶 Medium | Embeddings/Rebuild | 1-2 days | Vector database, embedding models |
| 🔶 Medium | Insights Generation | 1-2 days | Transformation system |
| 🔶 Lower | Context Management | 1 day | Context API endpoints |

**Total Estimated Effort: 12-18 days**

---

*Generated: [Current Date]*  
*Last Updated: [Current Date]*  
*Status: Initial Analysis Complete*