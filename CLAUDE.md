# Claude AI - Project Workflow Guide

## 🎯 Project Overview

**OpenNotebook CLI** - Go CLI-Anwendung zur Steuerung von OpenNotebook API (localhost:5055)
- Repo: `github.com/denkhaus/open-notebook-cli`
- Tech Stack: Go, urfav/cli/v2, samber/do/v2, zap-logger
- Architecture: HTTP Client mit DI für OpenNotebook API

---

## 📋 KNOT Project Management Workflow

### KNOT ist die **einzige Source of Truth** für alle Projektmanagement-Aufgaben.

### 🚀 **Mandatory Workflow Rules**

#### 1. **Immer KNOT verwenden**
- Alle Projekt-Tasks werden im KNOT-Tool verwaltet
- Keine Tasks ohne KNOT-Erfassung anfangen
- Jede Entscheidung im Projekt wird über KNOT dokumentiert

#### 2. **Dynamische Task-Erstellung**
- **Immer neue Tasks erstellen**, wenn neue Arbeiten entstehen
- Tasks müssen **sofort bei Erkennung** erstellt werden
- Beispiele: Bugfixes, neue Requirements, Architektur-Entscheidungen

#### 3. **Task Breakdown Pflicht**
- Complex Tasks (≥8) **müssen** in Sub-Tasks aufgeteilt werden
- `knot breakdown` - prüft, welche Tasks Breakdown brauchen
- Jede komplexe Arbeit wird in handhabbare Schritte zerlegt

#### 4. **Dependencies Pflicht**
- Tasks mit `knot dependency add` logisch verknüpfen
- Abhängigkeiten klar definieren und einhalten
- `knot actionable` - zeigt nächste verfügbare Tasks

#### 5. **State Management**
- Tasks durchgehen: `pending → in-progress → completed`
- **Immer** State aktualisieren bei Arbeitsbeginn/ende
- Keine Tasks ohne klaren State belassen

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
knot breakdown                    # Zeigt Tasks die Breakdown brauchen
knot actionable                   # Nächste verfügbare Tasks
knot blocked                      # Blockierte Tasks
```

### **Templates**
```bash
knot template list
knot template apply --name <template-name>
```

---

## 📊 **Project-Specific Workflow**

### **Task-Erstellung bei Bedarf**
**Immer neue Tasks erstellen, wenn:**
- ✅ Neue Requirements aus Kunden/Stakeholdern
- ✅ Architektur-Entscheidungen notwendig
- ✅ Tests für Features geschrieben werden müssen
- ✅ Documentation aktualisiert werden muss
- ✅ CI/CD Pipeline erweitert wird
- ✅ Performance Optimierungen nötig sind
- ✅ Security Issues gefunden werden
- ✅ Refactoring notwendig ist

**Beispiel-Nachricht:**
> "Während der Implementation von Feature X habe ich festgestellt, dass wir zusätzlich Y und Z implementieren müssen. Ich erstelle dafür separate Tasks."

### **Task Quality Standards**
- **Detaillierte Beschreibungen** mit API-Referenzen
- **Spezifische Complexity** (1-10)
- **Dependencies** klar definiert
- **Python-Referenzen** für API-Client Tasks
- **Sub-Tasks** für komplexe Aufgaben

### **Task Priorität & Complexität**
- **1-3**: Klein, schnell erledigt
- **4-6**: Mittel, Standard-Tasks
- **7-8**: Komplex, braucht Planning
- **9-10**: Sehr komplex, muss in Sub-Tasks broken down werden

---

## 🔄 **Typische Task-Sequenzen**

### **Neues Feature**
1. Research → Analysis → Design → Implementation → Testing → Documentation
2. Jeder Schritt = eigener Task mit Dependencies
3. Bei Komplexität → Sub-Tasks erstellen

### **Bug Fix**
1. Investigation → Root Cause → Fix Implementation → Testing → Verification
2. Immer separate Tasks für Debugging vs Fixing

### **Architecture Change**
1. Research → Proposal → Review → Implementation → Migration
2. Jede Phase = Task mit Dependencies

---

## 🎯 **Current Project Context**

**Laufendes Projekt:** OpenNotebook CLI (ID: `d27ada3e-7799-41bb-b7ed-200370663b5a`)

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
- [ ] **Jede Arbeit zuerst als Task erstellen**
- [ ] **Complex Tasks sofort breakdown**
- [ ] **Dependencies immer setzen**
- [ ] **State Management pflegen**
- [ ] **KNOT als Single Source of Truth**

### **NEVER:**
- [ ] Arbeiten ohne Task-Erfassung beginnen
- [ ] Dependencies ignorieren
- [ ] Komplexe Tasks ohne Breakdown lassen
- [ ] Task-States veralten lassen

---

## 📞 **Hilfe & Support**

- `knot get-started` - Umfassende Hilfe
- `knot help <command>` - Spezifische Hilfe
- Projekt-Dokumentation für Feature-Spezifika

**Remember:** KNOT ist nicht optional - KNOT ist unser zentrales Management-System!