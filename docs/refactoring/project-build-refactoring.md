# Project Build Refactoring Proposal

## Übersicht

Dieses Dokument enthält eine umfassende Analyse der `cmd/project-build.go` Datei und Vorschläge zur Verbesserung der Wartbarkeit und Erweiterbarkeit des Build-Prozesses.

**Analysiert am:** 2025-09-05  
**Datei:** `cmd/project-build.go`  
**Größe:** ~600+ Zeilen Code

## Aktuelle Probleme

### 1. Monolithische Funktionen
- **`RunE` Funktion (Zeilen 35-100)**: Zu viele Verantwortlichkeiten
  - Database connection
  - Project loading
  - Song querying
  - Configuration handling
  - Build orchestration

- **`buildProject` Funktion (Zeilen 102-200)**: Komplexe Orchestrierung
  - Directory creation
  - Song processing
  - Copyright handling
  - TOC creation
  - PDF merging

- **`buildSong` Funktion (Zeilen 300-450)**: Über 150 Zeilen
  - File reading
  - Config extraction
  - Merging
  - Temp file management
  - External tool execution
  - File operations

### 2. Globale Variablen (Zeilen 25-32)
```go
var (
    projectBuildOutputDir         string
    projectBuildAbcFileDir        string
    projectBuildPriorityThreshold int
    projectSampleId               string
)
```
**Probleme:**
- Erschweren Tests
- Können zu unerwarteten Seiteneffekten führen
- Machen parallele Ausführung schwierig

### 3. Hartcodierte Strukturen
- Verzeichnisnamen sind fest eingebaut
- Schwer anpassbare Ausgabestruktur
- Keine Flexibilität für verschiedene Projekt-Typen

### 4. Fehlende Abstraktion
- Keine klare Trennung zwischen Geschäftslogik und technischen Details
- PDF-Operationen sind über den Code verstreut
- Konfigurationshandling ist nicht zentralisiert

## Vorgeschlagene Refaktorierung

### 1. Builder Pattern

```go
type ProjectBuilder struct {
    config     *BuildConfig
    project    *ent.Project
    outputDirs *DirectoryStructure
    processor  *SongProcessor
    pdfManager *PDFManager
}

func NewProjectBuilder(config *BuildConfig) *ProjectBuilder {
    return &ProjectBuilder{
        config:     config,
        outputDirs: NewDirectoryStructure(config.OutputDir),
        processor:  NewSongProcessor(),
        pdfManager: NewPDFManager(),
    }
}

func (pb *ProjectBuilder) Build(ctx context.Context, projectID int) error {
    pipeline := NewBuildPipeline(
        &LoadProjectStep{},
        &SetupDirectoriesStep{},
        &ProcessSongsStep{},
        &CreateTOCStep{},
        &HandleCopyrightStep{},
        &MergePDFsStep{},
    )
    return pipeline.Execute(ctx, pb)
}
```

### 2. Konfiguration strukturieren

```go
type BuildConfig struct {
    OutputDir         string
    ABCFileDir        string
    PriorityThreshold int
    SampleID          string
    FolderPatterns    map[string]string
    Concurrency       int
}

type DirectoryStructure struct {
    Base         string
    PDF          string
    ABC          string
    Log          string
    Druckdateien string
    Referenz     string
}

func NewDirectoryStructure(baseDir string) *DirectoryStructure {
    return &DirectoryStructure{
        Base:         baseDir,
        PDF:          filepath.Join(baseDir, "pdf"),
        ABC:          filepath.Join(baseDir, "abc"),
        Log:          filepath.Join(baseDir, "log"),
        Druckdateien: filepath.Join(baseDir, "druckdateien"),
        Referenz:     filepath.Join(baseDir, "referenz"),
    }
}
```

### 3. Pipeline Pattern für Build-Schritte

```go
type BuildStep interface {
    Execute(ctx context.Context, builder *ProjectBuilder) error
    Name() string
}

type BuildPipeline struct {
    steps []BuildStep
}

func (bp *BuildPipeline) Execute(ctx context.Context, builder *ProjectBuilder) error {
    for _, step := range bp.steps {
        slog.Info("executing build step", "step", step.Name())
        if err := step.Execute(ctx, builder); err != nil {
            return fmt.Errorf("step %s failed: %w", step.Name(), err)
        }
    }
    return nil
}
```

### 4. Separate Services

#### SongProcessor
```go
type SongProcessor struct {
    configManager *ConfigurationManager
    zupfnoter     *ZupfnoterRunner
}

func (sp *SongProcessor) ProcessSong(ctx context.Context, song *ent.ProjectSong, config *BuildConfig) error {
    // Extrahiert aus buildSong Funktion
}
```

#### ConfigurationManager
```go
type ConfigurationManager struct {
    defaultConfig map[string]any
}

func (cm *ConfigurationManager) MergeConfigs(projectConfig, fileConfig map[string]any) (map[string]any, error) {
    // Konfigurationslogik
}
```

#### DirectoryManager
```go
type DirectoryManager struct {
    structure *DirectoryStructure
}

func (dm *DirectoryManager) SetupDirectories() error {
    // Verzeichniserstellung
}

func (dm *DirectoryManager) CleanupDirectories() error {
    // Aufräumen
}
```

## Implementierungsplan

### Phase 1: Grundstruktur
1. Erstelle neue Package-Struktur
2. Definiere Interfaces
3. Implementiere BuildConfig und DirectoryStructure

### Phase 2: Service Extraktion
1. Extrahiere PDFManager
2. Extrahiere SongProcessor
3. Extrahiere ConfigurationManager

### Phase 3: Pipeline Implementation
1. Implementiere BuildStep Interface
2. Erstelle konkrete Steps
3. Implementiere BuildPipeline

### Phase 4: Integration
1. Refaktoriere main command
2. Update Tests
3. Dokumentation

## Vorteile der Refaktorierung

✅ **Bessere Testbarkeit**: Kleinere, fokussierte Funktionen  
✅ **Erweiterbarkeit**: Neue Build-Schritte können einfach hinzugefügt werden  
✅ **Wartbarkeit**: Klare Verantwortlichkeiten und Abhängigkeiten  
✅ **Konfigurierbarkeit**: Flexible Anpassung ohne Code-Änderungen  
✅ **Parallelisierung**: Bessere Kontrolle über Nebenläufigkeit  
✅ **Fehlerbehandlung**: Strukturierte Behandlung von Fehlern  
✅ **Wiederverwendbarkeit**: Services können in anderen Kontexten genutzt werden

## Risiken und Überlegungen

⚠️ **Breaking Changes**: Bestehende Workflows könnten betroffen sein  
⚠️ **Komplexität**: Mehr Dateien und Strukturen  
⚠️ **Migration**: Bestehende Konfigurationen müssen angepasst werden  
⚠️ **Testing**: Umfangreiche Tests erforderlich

## Nächste Schritte

1. **Diskussion**: Team-Review der Vorschläge
2. **Prototyping**: Kleine Proof-of-Concept Implementation
3. **Schrittweise Migration**: Nicht alles auf einmal ändern
4. **Backward Compatibility**: Bestehende Funktionalität erhalten

---

*Dieses Dokument dient als Grundlage für zukünftige Refaktorierung-Entscheidungen und kann bei Bedarf erweitert werden.*