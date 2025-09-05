# Code Analysis: project-build.go

## Datei-Übersicht
- **Pfad:** `cmd/project-build.go`
- **Zeilen:** ~600+
- **Hauptfunktionen:** 15+
- **Komplexität:** Hoch

## Funktionsanalyse

### Hauptfunktionen und ihre Probleme

| Funktion | Zeilen | Verantwortlichkeiten | Probleme |
|----------|--------|---------------------|----------|
| `RunE` (Command) | 35-100 | DB-Verbindung, Projekt laden, Build orchestrieren | Zu viele Aufgaben |
| `buildProject` | 102-200 | Verzeichnisse, Songs, Copyright, TOC, PDF-Merge | Monolithisch |
| `buildSong` | 300-450 | ABC lesen, Config mergen, Zupfnoter, Dateien kopieren | 150+ Zeilen |
| `distributeZupfnoterOutput` | 380-450 | PDF-Verteilung nach Patterns | Pattern-Logik vermischt |
| `mergePDFs` | 480-520 | PDF-Zusammenführung | Fehlerbehandlung unvollständig |
| `copyPdfsToCopyrightDirectories` | 550-600 | Copyright-Organisation | Verschachtelte Logik |

### Globale Variablen (Anti-Pattern)
```go
var (
    projectBuildOutputDir         string  // Ausgabeverzeichnis
    projectBuildAbcFileDir        string  // ABC-Dateien Verzeichnis  
    projectBuildPriorityThreshold int     // Prioritätsschwelle
    projectSampleId               string  // Sample-ID
)
```

**Probleme:**
- Keine Kapselung
- Schwer testbar
- Race Conditions möglich
- Unklare Abhängigkeiten

### Hartcodierte Werte

#### Verzeichnisstruktur
```go
os.RemoveAll(filepath.Join(outputDir, "pdf"))
os.RemoveAll(filepath.Join(outputDir, "abc"))
os.RemoveAll(filepath.Join(outputDir, "log"))
os.RemoveAll(filepath.Join(outputDir, "druckdateien"))
os.RemoveAll(filepath.Join(outputDir, "referenz"))
```

#### Folder Patterns
```go
folderPatterns := map[string]string{
    "*_-A*_a3.pdf": "klein",
    "*_-M*_a3.pdf": "klein", 
    "*_-O*_a3.pdf": "klein",
    "*_-B*_a3.pdf": "gross",
    "*_-X*_a3.pdf": "gross",
}
```

## Abhängigkeitsanalyse

### Externe Abhängigkeiten
- `github.com/pdfcpu/pdfcpu/pkg/api` - PDF-Operationen
- `dario.cat/mergo` - Konfiguration mergen
- `golang.org/x/sync/errgroup` - Parallelisierung
- `github.com/bwl21/zupfmanager/internal/zupfnoter` - Externes Tool

### Interne Abhängigkeiten
- `internal/database` - Datenbankzugriff
- `internal/ent` - Entity Framework
- `github.com/spf13/cobra` - CLI Framework

## Datenfluss-Analyse

```
Command Input (Project ID)
    ↓
Database Query (Project + Songs)
    ↓
Directory Setup
    ↓
Parallel Song Processing
    ├── ABC File Reading
    ├── Config Extraction & Merging
    ├── Zupfnoter Execution
    ├── PDF Distribution
    └── File Copying
    ↓
Copyright Organization
    ↓
Table of Contents Creation
    ↓
PDF Merging by Patterns
    ↓
Cleanup
```

## Komplexitäts-Hotspots

### 1. buildSong Funktion
- **Zyklomatische Komplexität:** Hoch
- **Verantwortlichkeiten:** 8+
- **Abhängigkeiten:** 5+
- **Fehlerbehandlung:** Verstreut

### 2. distributeZupfnoterOutput
- **Pattern Matching:** Komplex
- **File Operations:** Viele
- **Error Handling:** Unvollständig

### 3. copyPdfsToCopyrightDirectories
- **Verschachtelte Schleifen:** 3 Ebenen
- **File Walking:** Ineffizient
- **Memory Usage:** Hoch bei vielen Dateien

## Performance-Probleme

### 1. File Operations
- Viele einzelne `os.MkdirAll` Aufrufe
- Ineffiziente `filepath.Walk` Nutzung
- Keine Batch-Operationen

### 2. Parallelisierung
- Nur auf Song-Ebene
- Keine PDF-Operation Parallelisierung
- Feste Limit von 5 Goroutines

### 3. Memory Usage
- Alle Songs werden in Memory geladen
- Große PDF-Dateien werden komplett gelesen
- Keine Streaming-Operationen

## Testbarkeit-Probleme

### 1. Globale State
- Globale Variablen erschweren Tests
- Keine Isolation zwischen Tests
- Setup/Teardown komplex

### 2. Externe Abhängigkeiten
- Dateisystem-Operationen
- Externe Tools (Zupfnoter)
- Datenbank-Zugriffe

### 3. Große Funktionen
- Schwer zu mocken
- Viele Pfade zu testen
- Setup aufwendig

## Wartbarkeits-Probleme

### 1. Code-Duplikation
- Directory creation logic
- Error handling patterns
- File operation patterns

### 2. Fehlende Abstraktion
- PDF-Operationen verstreut
- Konfiguration nicht zentralisiert
- Keine klaren Interfaces

### 3. Dokumentation
- Wenige Kommentare
- Keine API-Dokumentation
- Unklare Funktions-Contracts

## Erweiterbarkeits-Probleme

### 1. Neue Build-Schritte
- Müssen in monolithische Funktion eingefügt werden
- Keine Plugin-Architektur
- Schwer zu konfigurieren

### 2. Neue Output-Formate
- Hartcodierte PDF-Logik
- Keine Format-Abstraktion
- Schwer erweiterbar

### 3. Neue Konfigurationen
- Fest eingebaute Patterns
- Keine dynamische Konfiguration
- Schwer anpassbar

## Empfohlene Metriken für Refaktorierung

### Vor Refaktorierung
- **Funktionslänge:** 150+ Zeilen (buildSong)
- **Zyklomatische Komplexität:** >15
- **Abhängigkeiten:** >10 pro Funktion
- **Test Coverage:** Niedrig (geschätzt <30%)

### Nach Refaktorierung (Ziel)
- **Funktionslänge:** <50 Zeilen
- **Zyklomatische Komplexität:** <10
- **Abhängigkeiten:** <5 pro Funktion  
- **Test Coverage:** >80%

---

*Diese Analyse bildet die Grundlage für die Refaktorierung-Entscheidungen und sollte regelmäßig aktualisiert werden.*