# Claude AI - Project Workflow Guide

## 🎯 Project Overview

**OpenNotebook CLI** - Go CLI application for controlling OpenNotebook API (localhost:5055)

- Repo: `github.com/denkhaus/open-notebook-cli`
- Tech Stack: Go, urfav/cli/v2, samber/do/v2, zap-logger
- Architecture: HTTP Client with DI for OpenNotebook API
- All information about the API, like models, routes MUST be explored from [sourcecode](https://raw.githubusercontent.com/lfnovo/open-notebook/main/api)

---

## 🧠⚡ **INTEGRATED WORKFLOW SYSTEM**

### **Dual-Layer Architecture**
- **KNOT Tool**: **ONLY Source of Truth** for task management (Drop-in Replacement for Todo Tool)
- **Brain Tools**: Memory layer for knowledge management, user decisions, and project context

### **🚨 MANDATORY TOOL USAGE RULES**

#### 1. **KNOT ONLY - No Todo Tool Anymore**
- **KNOT is the official drop-in replacement for Claude Code's Todo Tool**
- **Todo Tool must NOT be used for any task management**
- All task operations MUST use KNOT commands exclusively
- KNOT provides superior hierarchical task management, state tracking, and auto-parent completion

#### 2. **Always Use KNOT First**
- All project tasks are managed in the KNOT tool
- Never start tasks without KNOT recording
- Every project decision is documented via KNOT
- Tasks go through: `pending → in-progress → completed`
- **Always** update state when starting/ending work

#### 3. **Dynamic Task Creation**
- **Always create new tasks** when new work arises
- Tasks must be created **immediately upon recognition**
- Examples: Bugfixes, new requirements, architecture decisions

#### 4. **Task Breakdown Requirement**
- Complex Tasks (≥8) **must** be broken down into sub-tasks
- `knot breakdown` - checks which tasks need breakdown
- Every complex work is broken down into manageable steps

#### 5. **Dependencies Requirement**
- Link tasks logically with `knot dependency add`
- Define and adhere to dependencies clearly
- `knot actionable` - shows next available tasks

---

## 🧠 **Brain Tools Integration**

### **Proactive Memory Storage Strategy**
**Critical**: Immediately store important information to brain tools when the user provides:
- Technical preferences (languages, tools, frameworks)
- Coding style or patterns
- Project requirements or constraints
- User opinions or feedback
- Problem-solving approaches
- Learning style or experience level
- Important decisions made during development
- Best practices for this project

### **Brain Tools Purpose**
The brain tools serve as a comprehensive knowledge database that:
- Stores user-provided information during programming process
- Maintains important decisions and principles across sessions
- Preserves project-specific best practices
- Provides context for future programming sessions
- Enables continuity of knowledge over time

### **Project Setup with Brain Tools**
Always load the project UUID from `.project` file when starting work and retrieve existing memories for context.

---

## 🔧 **Core KNOT Commands**

### **Project Management**

```bash
knot project create --title "Project Name" --description "Description"
knot project select --id <project-id>
knot project get-selected
```

### **Task Management**

```bash
knot task create --title "Task Title" --description "Details" --complexity 5
knot task list --depth-max 3
knot task update-state --id <task-id> --state in-progress|completed
knot task update-description --id <task-id> --description "New description"
```

### **Dependencies & Breakdown**

```bash
knot dependency add --task-id <task-id> --depends-on <dependency-id>
knot breakdown                    # Shows tasks that need breakdown
knot actionable                   # Next available tasks
knot blocked                      # Blocked tasks
```

### **Templates**

```bash
knot template list
knot template apply --name <template-name>
```

---

## 📊 **Project-Specific Workflow**

### **Task Creation as Needed**

**Always create new tasks when:**

- ✅ New requirements from customers/stakeholders
- ✅ Architecture decisions are necessary
- ✅ Tests need to be written for features
- ✅ Documentation needs to be updated
- ✅ CI/CD pipeline needs to be extended
- ✅ Performance optimizations are needed
- ✅ Security issues are found
- ✅ Refactoring is necessary

**Example Message:**

> "During the implementation of feature X, I discovered that we need to additionally implement Y and Z. I will create separate tasks for this."

### **Task Quality Standards**

- **Detailed descriptions** with API references
- **Specific complexity** (1-10)
- **Dependencies** clearly defined
- **Python references** for API client tasks
- **Sub-tasks** for complex tasks

### **Task Priority & Complexity**

- **1-3**: Small, quick to complete
- **4-6**: Medium, standard tasks
- **7-8**: Complex, needs planning
- **9-10**: Very complex, must be broken down into sub-tasks

---

## 🔄 **Typical Task Sequences**

### **New Feature**

1. Research → Analysis → Design → Implementation → Testing → Documentation
2. Each step = separate task with dependencies
3. If complex → create sub-tasks

### **Bug Fix**

1. Investigation → Root Cause → Fix Implementation → Testing → Verification
2. Always separate tasks for debugging vs fixing

### **Architecture Change**

1. Research → Proposal → Review → Implementation → Migration
2. Each phase = task with dependencies

---

## 🎯 **Current Project Context**

**Current Project:** OpenNotebook CLI (ID: `ab337d5c-078a-4d86-b186-3537e1c82947`)

**Completed:**

- ✅ API Research & Analysis
- ✅ Project Setup & Architecture

**Next Actions:**

- 🔄 CLI Command Structure Design (next actionable)
- 📋 API Client Development
- 📋 Authentication & Authorization

**Dependencies:**

- Project Setup → DI Setup → API Client → CLI Commands → Testing

---

## ⚠️ **Critical Rules**

### **MUST DO:**

- [ ] **Always create tasks first before any work**
- [ ] **Break down complex tasks immediately**
- [ ] **Always set dependencies**
- [ ] **Maintain state management**
- [ ] **Use KNOT as Single Source of Truth**

### **NEVER:**

- [ ] Start work without task recording
- [ ] Ignore dependencies
- [ ] Leave complex tasks without breakdown
- [ ] Let task states become outdated

---

## 📞 **Help & Support**

- `knot get-started` - Comprehensive help
- `knot help <command>` - Specific help
- Project documentation for feature specifics

---

## 🔄 **Cross-Reference System: Brain Tools ↔ KNOT Tasks**

### **Bidirectional Reference Capability**
**Critical Enhancement**: Brain memories and KNOT tasks/projects can reference each other for complete traceability:

#### **Memory → Task References**
- Store KNOT task IDs in memories when referencing specific work items
- Enables quick lookup of detailed task context when retrieving memories
- Provides complete audit trail from knowledge to implementation

#### **Task → Memory References**
- Include memory IDs in KNOT task descriptions for supporting documentation
- Creates comprehensive task context with supporting knowledge base
- Enables detailed task understanding through linked knowledge base

### **Reference Implementation Patterns**

**Memory Storage with Task References**:
```
Store memory with task ID in content/metadata:
"Related to KNOT task: [task-id]"
"Supporting documentation for: [task-title]"
```

**Task Management with Memory References**:
```
Update task description with memory references:
"See brain memory: [memory-id] for detailed analysis"
"Related context stored in memory: [memory-title]"
```

### **Benefits of Cross-Reference System**
1. **Complete Traceability**: From decision (memory) to implementation (task)
2. **Rich Task Context**: Tasks reference detailed knowledge and analysis
3. **Knowledge Integration**: Memories understand their implementation impact
4. **Audit Trail**: Full history of decisions and their execution
5. **Context Preservation**: Detailed knowledge available during task execution

---

## 💻 **Development Workflow Guidelines**

### **Critical Development Requirements**

#### 1. **Tool Usage Policy**
**KNOT as Todo Tool Replacement (MANDATORY)**:
- KNOT Tool is the **official drop-in replacement** for Todo Tool in Claude Code environments
- Todo Tool **must NOT be used** for task management anymore
- All task operations should use KNOT commands
- KNOT provides superior hierarchical task management, state tracking, and auto-parent completion

#### 2. **Cross-Reference Documentation System**
**Task ID and Memory ID References in Code**:
- Always include relevant Task IDs and Memory IDs in code comments
- This creates direct links between implementation and detailed documentation
- Enables quick navigation from code to task context and technical specifications
- Maintains traceability between implementation decisions and requirements

**Comment Examples**:
```go
// Auto-parent completion logic (Task ID: 06afc996-9a4e-4e75-a03d-8289d13042e3)
// See brain memory: 4bd7bc0a-4382-4fb3-8601-4facbbf1abc6 for technical specifications
func (s *service) evaluateAndUpdateParentTask(ctx context.Context, parentID uuid.UUID, actor string) error {
    // Implementation details...
}
```

#### 3. **Code Artifact Management**
**No Build Artifacts in Codebase**:
- All image operations, compilations, and build artifacts must use temporary directories
- Codebase must remain clean from compiled artifacts, images, and temporary files
- Use system temp directories or build-specific output directories

#### 4. **Implementation Quality Standards**
**Code Documentation Requirements**:
- Every significant function should reference related tasks/memories
- Technical decisions should link to brain memories with full specifications
- TODO comments should reference specific task IDs for follow-up
- This creates maintainable code with full context preservation

---

## 🧠 **Memory Maintenance Protocol**

### **Keeping Brain Up-to-Date**
**Critical**: Brain memories must be actively maintained to ensure accuracy:

1. **Memory Lifecycle Management**:
   - Regularly review and refine existing memories
   - Update memories when new information becomes available
   - Consolidate related memories into comprehensive entries
   - Delete outdated or obsolete memories

2. **Before Creating New Memories**:
   - Search for related existing memories first
   - Update existing memories instead of creating duplicates
   - Merge overlapping memories into coherent entries

3. **Memory Quality Standards**:
   - Keep memories concise but comprehensive
   - Update memories when project requirements change
   - Refine memories based on new user feedback

4. **Regular Maintenance**:
   - Audit memories for current relevance
   - Remove outdated information
   - Update memories when user preferences evolve
   - Use memory relationships to maintain consistency

This creates a growing knowledge base that ensures consistency and preserves the user's intent across all development sessions while preventing outdated memory accumulation.

---

**Remember:** KNOT is not optional - KNOT is our central management system! Brain Tools complement KNOT by providing the knowledge context needed for effective task execution.
