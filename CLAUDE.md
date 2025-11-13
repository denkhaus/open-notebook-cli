# Claude AI - Project Workflow Guide

## 🎯 Project Overview

**OpenNotebook CLI** - Go CLI application for controlling OpenNotebook API (localhost:5055)

- Repo: `github.com/denkhaus/open-notebook-cli`
- Tech Stack: Go, urfav/cli/v2, samber/do/v2, zap-logger
- Architecture: HTTP Client with DI for OpenNotebook API
- All information about the API, like models, routes MUST be explored from [sourcecode](https://raw.githubusercontent.com/lfnovo/open-notebook/main/api)

---

## 📋 KNOT Project Management Workflow

### KNOT is the **only Source of Truth** for all project management tasks.

### 🚀 **Mandatory Workflow Rules**

#### 1. **Always Use KNOT**

- All project tasks are managed in the KNOT tool
- Never start tasks without KNOT recording
- Every project decision is documented via KNOT

#### 2. **Dynamic Task Creation**

- **Always create new tasks** when new work arises
- Tasks must be created **immediately upon recognition**
- Examples: Bugfixes, new requirements, architecture decisions

#### 3. **Task Breakdown Requirement**

- Complex Tasks (≥8) **must** be broken down into sub-tasks
- `knot breakdown` - checks which tasks need breakdown
- Every complex work is broken down into manageable steps

#### 4. **Dependencies Requirement**

- Link tasks logically with `knot dependency add`
- Define and adhere to dependencies clearly
- `knot actionable` - shows next available tasks

#### 5. **State Management**

- Tasks go through: `pending → in-progress → completed`
- **Always** update state when starting/ending work
- Never leave tasks without clear state

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

**Current Project:** OpenNotebook CLI (ID: `d27ada3e-7799-41bb-b7ed-200370663b5a`)

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

**Remember:** KNOT is not optional - KNOT is our central management system!
