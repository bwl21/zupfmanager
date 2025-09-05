# PDFManager Specification

## Übersicht

Der **PDFManager** ist eine vorgeschlagene Komponente zur Kapselung aller PDF-bezogenen Operationen im Build-Prozess. Diese Spezifikation definiert die genauen Funktionen und Verantwortlichkeiten.

## Aktuelle PDF-Operationen (Ist-Zustand)

### In project-build.go identifizierte PDF-Operationen:
1. **distributeZupfnoterOutput** - Verteilung von generierten PDFs
2. **mergePDFs** - Zusammenführung von PDF-Dateien
3. **copyFile** - Kopieren einzelner PDF-Dateien
4. **copyPdfsToCopyrightDirectories** - Copyright-basierte Organisation
5. **Pattern Matching** - Klassifizierung nach Dateinamen-Mustern

## PDFManager Interface

### Core Interface
```go
type PDFManager interface {
    // Distribution & Organization
    DistributeFiles(ctx context.Context, request *DistributionRequest) error
    OrganizeByCopyright(ctx context.Context, request *CopyrightRequest) error
    
    // Merging Operations
    MergeByPattern(ctx context.Context, request *MergeRequest) error
    MergeBatch(ctx context.Context, requests []*MergeRequest) error
    
    // File Operations
    CopyFile(src, dst string) error
    ValidateFile(path string) error
    
    // Discovery & Pattern Matching
    FindFiles(pattern string, baseDir string) ([]string, error)
    MatchPattern(filename string, patterns map[string]string) (string, bool)
    
    // Cleanup & Maintenance
    CleanupTempFiles() error
    ValidateIntegrity(files []string) error
}
```

## Detaillierte Funktionsspezifikationen

### 1. File Distribution Functions

#### DistributeFiles
```go
func (pm *PDFManager) DistributeFiles(ctx context.Context, req *DistributionRequest) error
```
**Zweck:** Hauptfunktion für PDF-Verteilung nach Zupfnoter-Generierung  
**Input:** DistributionRequest mit Source-Dir, Patterns, Ziel-Verzeichnissen  
**Output:** Error bei Fehlern  
**Ersetzt:** `distributeZupfnoterOutput` Funktion

**Workflow:**
1. Finde alle PDFs im Source-Verzeichnis
2. Wende Folder-Patterns an
3. Benenne Dateien mit Index um
4. Kopiere in entsprechende Ziel-Verzeichnisse
5. Validiere Kopier-Operationen

#### OrganizeByFolderPatterns
```go
func (pm *PDFManager) OrganizeByFolderPatterns(files []string, patterns map[string]string) error
```
**Zweck:** Organisiert PDFs nach konfigurierbaren Mustern  
**Input:** Liste von PDF-Dateien, Pattern-Map  
**Output:** Error bei Pattern-Fehlern

**Standard-Patterns:**
- `*_-A*_a3.pdf` → "klein"
- `*_-M*_a3.pdf` → "klein"  
- `*_-O*_a3.pdf` → "klein"
- `*_-B*_a3.pdf` → "gross"
- `*_-X*_a3.pdf` → "gross"

### 2. Pattern Matching & Classification

#### ClassifyPDF
```go
func (pm *PDFManager) ClassifyPDF(filename string) PDFType
```
**Zweck:** Erkennt PDF-Typ basierend auf Dateinamen  
**Input:** Dateiname  
**Output:** PDFType (enum: Klein, Gross, TOC, Unknown)

#### ApplyFolderPatterns
```go
func (pm *PDFManager) ApplyFolderPatterns(filename string, patterns map[string]string) string
```
**Zweck:** Wendet Folder-Patterns auf Dateinamen an  
**Input:** Dateiname, Pattern-Map  
**Output:** Ziel-Ordner oder leerer String

#### FindSongPDFs
```go
func (pm *PDFManager) FindSongPDFs(songBasename string, searchDir string) ([]string, error)
```
**Zweck:** Findet alle PDFs die zu einem Song gehören  
**Input:** Song-Basename (ohne .abc), Such-Verzeichnis  
**Output:** Liste von PDF-Pfaden

### 3. Merging Operations

#### MergeDirectoryPDFs
```go
func (pm *PDFManager) MergeDirectoryPDFs(sourceDir, destFile string) error
```
**Zweck:** Führt alle PDFs in einem Verzeichnis zusammen  
**Input:** Quell-Verzeichnis, Ziel-Datei  
**Output:** Error bei Merge-Fehlern  
**Ersetzt:** `mergePDFs` Funktion

**Workflow:**
1. Scanne Verzeichnis nach PDF-Dateien
2. Sortiere Dateien alphabetisch
3. Validiere PDF-Integrität
4. Führe mit pdfcpu zusammen
5. Validiere Ergebnis-PDF

#### MergeFileList
```go
func (pm *PDFManager) MergeFileList(files []string, destFile string) error
```
**Zweck:** Führt spezifische PDF-Liste zusammen  
**Input:** Liste von PDF-Pfaden, Ziel-Datei  
**Output:** Error bei Merge-Fehlern

#### MergeBatchDirectories
```go
func (pm *PDFManager) MergeBatchDirectories(dirs map[string]string) error
```
**Zweck:** Batch-Merge für mehrere Verzeichnisse  
**Input:** Map von Source-Dir zu Dest-File  
**Output:** Error bei Batch-Fehlern

### 4. Copyright Management

#### OrganizeByCopyright
```go
func (pm *PDFManager) OrganizeByCopyright(project *ent.Project, outputDir string) error
```
**Zweck:** Kopiert PDFs in Copyright-spezifische Verzeichnisse  
**Input:** Projekt-Entity, Output-Verzeichnis  
**Output:** Error bei Copyright-Organisation  
**Ersetzt:** `copyPdfsToCopyrightDirectories` Funktion

**Workflow:**
1. Extrahiere Copyright-Informationen aus Songs
2. Erstelle Copyright-Verzeichnisse
3. Finde zugehörige PDF-Dateien
4. Kopiere in entsprechende Verzeichnisse
5. Validiere Kopier-Operationen

#### SetupCopyrightDirectories
```go
func (pm *PDFManager) SetupCopyrightDirectories(copyrights []string, baseDir string) error
```
**Zweck:** Erstellt Copyright-Verzeichnisstruktur  
**Input:** Liste von Copyright-Namen, Basis-Verzeichnis  
**Output:** Error bei Verzeichnis-Erstellung

### 5. File Operations & Utilities

#### CopyFile
```go
func (pm *PDFManager) CopyFile(src, dst string) error
```
**Zweck:** Sichere Datei-Kopie mit Validierung  
**Input:** Quell-Pfad, Ziel-Pfad  
**Output:** Error bei Kopier-Fehlern  
**Ersetzt:** `copyFile` Funktion

**Features:**
- Validierung der Quell-Datei
- Atomic Copy-Operation
- Permissions-Erhaltung
- Fehlerbehandlung

#### ValidatePDF
```go
func (pm *PDFManager) ValidatePDF(path string) error
```
**Zweck:** Validiert PDF-Datei-Integrität  
**Input:** PDF-Pfad  
**Output:** Error bei Validierungs-Fehlern

#### EnsureDirectory
```go
func (pm *PDFManager) EnsureDirectory(path string) error
```
**Zweck:** Erstellt Verzeichnisstruktur falls nötig  
**Input:** Verzeichnis-Pfad  
**Output:** Error bei Erstellung

### 6. Discovery & Search

#### FindPDFsByPattern
```go
func (pm *PDFManager) FindPDFsByPattern(baseDir, pattern string) ([]string, error)
```
**Zweck:** Findet alle PDFs in Verzeichnis mit Pattern  
**Input:** Basis-Verzeichnis, Such-Pattern  
**Output:** Liste von PDF-Pfaden, Error

#### GetPDFStats
```go
func (pm *PDFManager) GetPDFStats(dir string) (*PDFStats, error)
```
**Zweck:** Sammelt PDF-Statistiken  
**Input:** Verzeichnis-Pfad  
**Output:** PDFStats-Struktur, Error

#### ListAllPDFs
```go
func (pm *PDFManager) ListAllPDFs(baseDir string) ([]string, error)
```
**Zweck:** Listet alle PDF-Dateien rekursiv auf  
**Input:** Basis-Verzeichnis  
**Output:** Liste aller PDF-Pfade

### 7. Configuration & Patterns

#### LoadFolderPatterns
```go
func (pm *PDFManager) LoadFolderPatterns(project *ent.Project) map[string]string
```
**Zweck:** Lädt Folder-Patterns aus Projekt-Config  
**Input:** Projekt-Entity  
**Output:** Pattern-Map

#### ValidatePatterns
```go
func (pm *PDFManager) ValidatePatterns(patterns map[string]string) error
```
**Zweck:** Validiert Pattern-Konfiguration  
**Input:** Pattern-Map  
**Output:** Error bei ungültigen Patterns

#### GetDefaultPatterns
```go
func (pm *PDFManager) GetDefaultPatterns() map[string]string
```
**Zweck:** Liefert Standard-Patterns falls keine Config vorhanden  
**Output:** Standard-Pattern-Map

## Datenstrukturen

### DistributionRequest
```go
type DistributionRequest struct {
    SourceDir      string            // Verzeichnis mit generierten PDFs
    BaseFilename   string            // Basis-Dateiname (ohne Extension)
    OutputDir      string            // Haupt-Output-Verzeichnis
    SongIndex      int               // Index für Datei-Nummerierung
    FolderPatterns map[string]string // Pattern zu Ordner Mapping
}
```

### CopyrightRequest
```go
type CopyrightRequest struct {
    Project   *ent.Project      // Projekt-Entity
    OutputDir string            // Output-Verzeichnis
    Songs     []*ent.ProjectSong // Liste der Songs
}
```

### MergeRequest
```go
type MergeRequest struct {
    SourceDir string // Quell-Verzeichnis
    DestFile  string // Ziel-Datei
    Pattern   string // Optional: Filter-Pattern
}
```

### PDFStats
```go
type PDFStats struct {
    TotalFiles int                // Gesamtanzahl PDF-Dateien
    TotalSize  int64              // Gesamtgröße in Bytes
    ByType     map[string]int     // Anzahl nach Typ
    ByFolder   map[string]int     // Anzahl nach Ordner
    Errors     []string           // Liste von Fehlern
}
```

### PDFType
```go
type PDFType int

const (
    PDFTypeUnknown PDFType = iota
    PDFTypeKlein           // Kleine Formate (A, M, O)
    PDFTypeGross           // Große Formate (B, X)
    PDFTypeTOC             // Table of Contents
    PDFTypeReference       // Referenz-Dateien
)
```

## Error Handling

### Strukturierte Fehlerbehandlung
```go
type PDFError struct {
    Operation string // Operation die fehlschlug
    File      string // Betroffene Datei
    Cause     error  // Ursprünglicher Fehler
}

func (e *PDFError) Error() string {
    return fmt.Sprintf("PDF operation '%s' failed for file '%s': %v", 
        e.Operation, e.File, e.Cause)
}
```

### Logging
```go
func (pm *PDFManager) LogOperation(operation string, details map[string]interface{}) {
    slog.Info("PDF operation", 
        "operation", operation,
        "details", details)
}
```

## Performance Considerations

### Parallelisierung
- Batch-Operationen für bessere Performance
- Parallele PDF-Verarbeitung wo möglich
- Streaming für große Dateien

### Memory Management
- Lazy Loading von PDF-Metadaten
- Streaming-Kopie für große Dateien
- Cleanup von temporären Dateien

### Caching
- Pattern-Matching-Ergebnisse cachen
- PDF-Metadaten cachen
- Verzeichnis-Scans cachen

## Testing Strategy

### Unit Tests
- Jede Funktion isoliert testbar
- Mock-Filesystem für Tests
- Pattern-Matching Tests

### Integration Tests
- End-to-End PDF-Workflows
- Reale Datei-Operationen
- Performance-Tests

### Error Cases
- Ungültige PDF-Dateien
- Fehlende Verzeichnisse
- Permissions-Probleme

## Migration Path

### Phase 1: Interface Definition
- Definiere PDFManager Interface
- Erstelle Basis-Implementierung
- Schreibe Tests

### Phase 2: Funktions-Extraktion
- Extrahiere PDF-Operationen aus project-build.go
- Implementiere PDFManager-Funktionen
- Update Aufrufer

### Phase 3: Optimierung
- Performance-Verbesserungen
- Erweiterte Fehlerbehandlung
- Zusätzliche Features

---

*Diese Spezifikation dient als Blaupause für die PDFManager-Implementierung und kann bei Bedarf erweitert werden.*