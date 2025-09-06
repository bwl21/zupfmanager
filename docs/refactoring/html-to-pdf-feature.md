# HTML to PDF Conversion Feature

## Anforderung

Parallel zur ABC-zu-PDF Konvertierung mit Zupfnoter soll eine entsprechende HTML-Datei (gleicher Dateiname wie ABC im Source-Ordner) zu PDF konvertiert werden.

### Spezifikationen:
- **HTML-Quelle:** Von Zupfnoter generierte HTML-Dateien (gleicher Dateiname wie ABC-Datei)
- **Blattnummer:** Dynamische Injection via DOM-Manipulation (nicht in Original-HTML)
- **Original-Schutz:** HTML-Dateien dürfen NICHT verändert werden (Zupfnoter-generiert)
- **DOM-Injection:** CSS + HTML-Element wird vor PDF-Generierung in `<body>` eingefügt
- **Konvertierung:** Mit Go-Paket "chromedp"
- **Parallelisierung:** Parallel zur ABC-Konvertierung
- **Erweiterbarkeit:** Weitere DOM-Manipulationen sind zu erwarten

## Architektur-Vorschlag

### 1. HTML-zu-PDF Converter Interface

```go
package htmlpdf

import (
    "context"
    "github.com/bwl21/zupfmanager/internal/ent"
)

type HTMLToPDFConverter interface {
    ConvertToPDF(ctx context.Context, request *ConversionRequest) (*ConversionResult, error)
    ValidateHTML(htmlPath string) error
    Close() error
}

type ConversionRequest struct {
    HTMLFilePath  string            // Pfad zur HTML-Quelldatei (Zupfnoter-generiert)
    OutputPath    string            // Ziel-PDF-Pfad
    SongIndex     int               // Blattnummer für DOM-Injection
    Song          *ent.ProjectSong  // Song-Informationen
    DOMInjectors  []DOMInjector     // Liste von DOM-Injektoren
    DOMScripts    []string          // JavaScript-Code für DOM-Manipulation
}

type ConversionResult struct {
    OutputPath   string        // Pfad zur generierten PDF
    PageCount    int           // Anzahl Seiten
    FileSize     int64         // Dateigröße
    Duration     time.Duration // Konvertierungszeit
    Warnings     []string      // Warnungen
}

type HTMLTransform interface {
    Apply(htmlContent string, request *ConversionRequest) (string, error)
    Name() string
}
```

### 2. ChromeDP Implementation

```go
package htmlpdf

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "time"
    
    "github.com/chromedp/chromedp"
)

type ChromeDPConverter struct {
    allocCtx   context.Context
    cancelFunc context.CancelFunc
    injectors  []DOMInjector
}

func NewChromeDPConverter(injectors ...DOMInjector) *ChromeDPConverter {
    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
        chromedp.NoSandbox,
        chromedp.Headless,
        chromedp.DisableGPU,
        chromedp.NoDefaultBrowserCheck,
        chromedp.Flag("disable-background-timer-throttling", true),
        chromedp.Flag("disable-backgrounding-occluded-windows", true),
        chromedp.Flag("disable-renderer-backgrounding", true),
        chromedp.Flag("disable-web-security", true), // Für lokale Dateien
    )
    
    return &ChromeDPConverter{
        allocCtx:   allocCtx,
        cancelFunc: cancel,
        injectors:  injectors,
    }
}

func (c *ChromeDPConverter) ConvertToPDF(ctx context.Context, request *ConversionRequest) (*ConversionResult, error) {
    start := time.Now()
    
    // 1. Validiere dass HTML-Datei existiert (Zupfnoter-generiert)
    if _, err := os.Stat(request.HTMLFilePath); os.IsNotExist(err) {
        return nil, fmt.Errorf("HTML file does not exist: %s", request.HTMLFilePath)
    }
    
    // 2. Bereite DOM-Injektoren vor
    request.DOMScripts = make([]string, 0)
    for _, injector := range c.injectors {
        err := injector.InjectIntoDOM(ctx, request)
        if err != nil {
            return nil, fmt.Errorf("DOM injector %s failed: %w", injector.Name(), err)
        }
    }
    
    // 3. PDF generieren mit DOM-Manipulation
    var pdfBuffer []byte
    taskCtx, cancel := chromedp.NewContext(c.allocCtx)
    defer cancel()
    
    // Erstelle ChromeDP-Aktionen
    actions := []chromedp.Action{
        chromedp.Navigate("file://" + request.HTMLFilePath),
        chromedp.WaitReady("body"),
    }
    
    // Füge DOM-Manipulations-Skripte hinzu
    for _, script := range request.DOMScripts {
        actions = append(actions, chromedp.Evaluate(script, nil))
    }
    
    // Warte kurz damit DOM-Änderungen wirksam werden
    actions = append(actions, chromedp.Sleep(500*time.Millisecond))
    
    // PDF-Generierung
    actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
        var err error
        pdfBuffer, _, err = chromedp.PrintToPDF().
            WithPrintBackground(true).
            WithPaperWidth(8.27).    // A4 width in inches
            WithPaperHeight(11.7).   // A4 height in inches
            WithMarginTop(0.4).
            WithMarginBottom(0.4).
            WithMarginLeft(0.4).
            WithMarginRight(0.4).
            WithDisplayHeaderFooter(false).
            Do(ctx)
        return err
    }))
    
    err := chromedp.Run(taskCtx, actions...)
    if err != nil {
        return nil, fmt.Errorf("failed to generate PDF: %w", err)
    }
    
    // 4. PDF-Datei schreiben
    err = os.WriteFile(request.OutputPath, pdfBuffer, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to write PDF: %w", err)
    }
    
    return &ConversionResult{
        OutputPath: request.OutputPath,
        FileSize:   int64(len(pdfBuffer)),
        Duration:   time.Since(start),
    }, nil
}

// createTempHTML ist nicht mehr nötig, da wir die Original-HTML-Datei direkt verwenden
// und DOM-Manipulation zur Laufzeit durchführen

func (c *ChromeDPConverter) Close() error {
    if c.cancelFunc != nil {
        c.cancelFunc()
    }
    return nil
}
```

### 3. DOM-Injection für Zupfnoter HTML-Dateien

Da die HTML-Dateien von Zupfnoter generiert werden und nicht verändert werden dürfen, erfolgt die Blattnummer-Injection über DOM-Manipulation zur Laufzeit.

#### DOM-Injection Interface
```go
package htmlpdf

import (
    "context"
    "fmt"
)

type DOMInjector interface {
    InjectIntoDOM(ctx context.Context, request *ConversionRequest) error
    Name() string
}

type PageNumberInjector struct {
    cssStyle string
    position string // "top-right", "bottom-right", etc.
}

func NewPageNumberInjector(position string) *PageNumberInjector {
    cssStyle := `
        <style>
        @media print {
            #druckParagraph {
                position: fixed;
                bottom: 0;
                right: 0;
                margin: 10px;
                font-weight: bold;
                background: grey;
                padding: 5px;
                border: 1px solid black;
                z-index: 1000;
            }
        }
        </style>
    `
    
    return &PageNumberInjector{
        cssStyle: cssStyle,
        position: position,
    }
}

func (inj *PageNumberInjector) InjectIntoDOM(ctx context.Context, request *ConversionRequest) error {
    // Diese Funktion wird von ChromeDP aufgerufen
    pageNumber := fmt.Sprintf("%02d", request.SongIndex)
    
    // 1. Entferne <text> Elemente mit '#vb' Inhalt
    removeVbScript := `
        const textElements = document.querySelectorAll('text');
        textElements.forEach(element => {
            if (element.textContent && element.textContent.trim() === '#vb') {
                element.remove();
            }
        });
    `
    
    // 2. CSS-Style in <head> einfügen
    styleScript := fmt.Sprintf(`
        const style = document.createElement('style');
        style.textContent = %s;
        document.head.appendChild(style);
    `, "`"+inj.cssStyle+"`")
    
    // 3. HTML-Element am Anfang von <body> einfügen
    elementScript := fmt.Sprintf(`
        const paragraph = document.createElement('p');
        paragraph.id = 'druckParagraph';
        paragraph.textContent = '%s';
        document.body.insertBefore(paragraph, document.body.firstChild);
    `, pageNumber)
    
    // Alle Skripte werden von ChromeDP ausgeführt
    request.DOMScripts = append(request.DOMScripts, removeVbScript, styleScript, elementScript)
    
    return nil
}

func (inj *PageNumberInjector) Name() string {
    return "PageNumberInjector"
}

// Spezifischer Injector für das Entfernen von '#vb' Text-Elementen
type TextCleanupInjector struct {
    removePatterns []string
}

func NewTextCleanupInjector(patterns ...string) *TextCleanupInjector {
    if len(patterns) == 0 {
        patterns = []string{"#vb"} // Standard-Pattern
    }
    return &TextCleanupInjector{
        removePatterns: patterns,
    }
}

func (inj *TextCleanupInjector) InjectIntoDOM(ctx context.Context, request *ConversionRequest) error {
    for _, pattern := range inj.removePatterns {
        removeScript := fmt.Sprintf(`
            const textElements = document.querySelectorAll('text');
            textElements.forEach(element => {
                if (element.textContent && element.textContent.trim() === '%s') {
                    element.remove();
                }
            });
        `, pattern)
        
        request.DOMScripts = append(request.DOMScripts, removeScript)
    }
    return nil
}

func (inj *TextCleanupInjector) Name() string {
    return "TextCleanupInjector"
}

// Erweiterte DOM-Injection für komplexere Anpassungen
type CustomDOMInjector struct {
    cssStyles     []string
    htmlElements  []HTMLElement
    cleanupRules  []CleanupRule
}

type CleanupRule struct {
    Selector string // CSS-Selector für Elemente
    Action   string // "remove", "hide", "modify"
    Pattern  string // Text-Pattern für Matching
    Value    string // Neuer Wert bei "modify"
}

type HTMLElement struct {
    Tag        string
    ID         string
    Class      string
    Content    string
    Attributes map[string]string
    Position   InsertPosition
}

type InsertPosition int

const (
    BodyStart InsertPosition = iota
    BodyEnd
    HeadEnd
    BeforeElement
    AfterElement
)

func NewCustomDOMInjector() *CustomDOMInjector {
    return &CustomDOMInjector{
        cssStyles:    make([]string, 0),
        htmlElements: make([]HTMLElement, 0),
        cleanupRules: make([]CleanupRule, 0),
    }
}

func (inj *CustomDOMInjector) AddCSS(css string) *CustomDOMInjector {
    inj.cssStyles = append(inj.cssStyles, css)
    return inj
}

func (inj *CustomDOMInjector) AddElement(element HTMLElement) *CustomDOMInjector {
    inj.htmlElements = append(inj.htmlElements, element)
    return inj
}

func (inj *CustomDOMInjector) AddCleanupRule(rule CleanupRule) *CustomDOMInjector {
    inj.cleanupRules = append(inj.cleanupRules, rule)
    return inj
}

func (inj *CustomDOMInjector) InjectIntoDOM(ctx context.Context, request *ConversionRequest) error {
    // 1. Cleanup-Regeln anwenden (zuerst)
    for _, rule := range inj.cleanupRules {
        cleanupScript := inj.generateCleanupScript(rule)
        request.DOMScripts = append(request.DOMScripts, cleanupScript)
    }
    
    // 2. CSS-Styles hinzufügen
    for _, css := range inj.cssStyles {
        styleScript := fmt.Sprintf(`
            const style = document.createElement('style');
            style.textContent = %s;
            document.head.appendChild(style);
        `, "`"+css+"`")
        request.DOMScripts = append(request.DOMScripts, styleScript)
    }
    
    // 3. HTML-Elemente hinzufügen
    for _, element := range inj.htmlElements {
        elementScript := inj.generateElementScript(element, request)
        request.DOMScripts = append(request.DOMScripts, elementScript)
    }
    
    return nil
}

func (inj *CustomDOMInjector) generateCleanupScript(rule CleanupRule) string {
    switch rule.Action {
    case "remove":
        if rule.Pattern != "" {
            return fmt.Sprintf(`
                const elements = document.querySelectorAll('%s');
                elements.forEach(element => {
                    if (element.textContent && element.textContent.trim() === '%s') {
                        element.remove();
                    }
                });
            `, rule.Selector, rule.Pattern)
        } else {
            return fmt.Sprintf(`
                const elements = document.querySelectorAll('%s');
                elements.forEach(element => element.remove());
            `, rule.Selector)
        }
    case "hide":
        return fmt.Sprintf(`
            const elements = document.querySelectorAll('%s');
            elements.forEach(element => {
                if (!element.textContent || element.textContent.trim() === '%s') {
                    element.style.display = 'none';
                }
            });
        `, rule.Selector, rule.Pattern)
    case "modify":
        return fmt.Sprintf(`
            const elements = document.querySelectorAll('%s');
            elements.forEach(element => {
                if (element.textContent && element.textContent.trim() === '%s') {
                    element.textContent = '%s';
                }
            });
        `, rule.Selector, rule.Pattern, rule.Value)
    default:
        return ""
    }
}

func (inj *CustomDOMInjector) generateElementScript(element HTMLElement, request *ConversionRequest) string {
    // Platzhalter ersetzen
    content := element.Content
    content = strings.ReplaceAll(content, "${SONG_INDEX}", fmt.Sprintf("%02d", request.SongIndex))
    content = strings.ReplaceAll(content, "${SONG_TITLE}", request.Song.Edges.Song.Title)
    
    script := fmt.Sprintf(`
        const element = document.createElement('%s');
        element.id = '%s';
        element.className = '%s';
        element.textContent = '%s';
    `, element.Tag, element.ID, element.Class, content)
    
    // Attribute hinzufügen
    for key, value := range element.Attributes {
        script += fmt.Sprintf(`element.setAttribute('%s', '%s');`, key, value)
    }
    
    // Position bestimmen
    switch element.Position {
    case BodyStart:
        script += `document.body.insertBefore(element, document.body.firstChild);`
    case BodyEnd:
        script += `document.body.appendChild(element);`
    case HeadEnd:
        script += `document.head.appendChild(element);`
    }
    
    return script
}

func (inj *CustomDOMInjector) Name() string {
    return "CustomDOMInjector"
}
```

### 4. Integration in buildSong Funktion

```go
// Erweiterte buildSong Funktion mit HTML-zu-PDF Unterstützung
func buildSong(ctx context.Context, abcFileDir, outputDir string, songIndex int, song *ent.ProjectSong, projectSampleId string, project *ent.Project) error {
    slog.Info("building song", "song", song.Edges.Song.Title)
    
    // Erstelle errgroup für parallele Verarbeitung
    eg, egCtx := errgroup.WithContext(ctx)
    
    // ABC-zu-PDF Konvertierung (bestehend)
    eg.Go(func() error {
        return buildSongABC(egCtx, abcFileDir, outputDir, songIndex, song, projectSampleId, project)
    })
    
    // HTML-zu-PDF Konvertierung (neu)
    eg.Go(func() error {
        return buildSongHTML(egCtx, abcFileDir, outputDir, songIndex, song, project)
    })
    
    // Warte auf beide Konvertierungen
    if err := eg.Wait(); err != nil {
        return fmt.Errorf("failed to build song %s: %w", song.Edges.Song.Title, err)
    }
    
    return nil
}

// Extrahierte ABC-Konvertierung
func buildSongABC(ctx context.Context, abcFileDir, outputDir string, songIndex int, song *ent.ProjectSong, projectSampleId string, project *ent.Project) error {
    // Bestehende ABC-zu-PDF Logik hier...
    // (Code aus der aktuellen buildSong Funktion)
    return nil
}

// Neue HTML-Konvertierung
func buildSongHTML(ctx context.Context, abcFileDir, outputDir string, songIndex int, song *ent.ProjectSong, project *ent.Project) error {
    // 1. Prüfe ob HTML-Datei existiert
    htmlFilename := strings.TrimSuffix(song.Edges.Song.Filename, ".abc") + ".html"
    htmlPath := filepath.Join(abcFileDir, htmlFilename)
    
    if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
        slog.Debug("no HTML file found for song", "song", song.Edges.Song.Title, "expected", htmlFilename)
        return nil // Kein Fehler, HTML ist optional
    }
    
    // 2. HTML-zu-PDF Converter mit DOM-Injektoren erstellen
    converter := NewChromeDPConverter(
        NewTextCleanupInjector("#vb"),           // Entfernt <text>#vb</text> Elemente
        NewPageNumberInjector("bottom-right"),   // Fügt Blattnummer hinzu
        // Weitere DOM-Injektoren können hier hinzugefügt werden
    )
    defer converter.Close()
    
    // 3. Output-Pfad bestimmen (gleicher Name wie ABC-Datei mit '_noten.pdf' Endung)
    abcBasename := strings.TrimSuffix(song.Edges.Song.Filename, ".abc")
    pdfFilename := abcBasename + "_noten.pdf"
    outputPath := filepath.Join(outputDir, "pdf", pdfFilename)
    
    // 4. Konvertierung durchführen (Original-HTML bleibt unverändert)
    request := &ConversionRequest{
        HTMLFilePath: htmlPath,    // Direkt von Zupfnoter generierte HTML-Datei
        OutputPath:   outputPath,
        SongIndex:    songIndex,
        Song:         song,
        DOMInjectors: converter.injectors, // Wird automatisch gesetzt
    }
    
    result, err := converter.ConvertToPDF(ctx, request)
    if err != nil {
        return fmt.Errorf("failed to convert HTML to PDF for %s: %w", song.Edges.Song.Title, err)
    }
    
    slog.Info("HTML to PDF conversion completed", 
        "song", song.Edges.Song.Title,
        "output", result.OutputPath,
        "filename", pdfFilename,
        "duration", result.Duration,
        "size", result.FileSize)
    
    // 5. PDF in Druckdateien-Verzeichnisse verteilen (analog zu ABC-PDFs)
    return distributeHTMLPDF(project, htmlFilename, outputDir, songIndex)
}
```

### 5. PDF-Verteilung für HTML-PDFs

```go
func distributeHTMLPDF(project *ent.Project, htmlFilename string, outputDir string, songIndex int) error {
    pdfDir := filepath.Join(outputDir, "pdf")
    // HTML-Dateiname entspricht ABC-Dateiname, daher verwenden wir den gleichen Basename
    baseFilenameWithoutExt := strings.TrimSuffix(htmlFilename, ".html")
    
    // Suche nach HTML-generierten PDFs mit '_noten.pdf' Endung
    pattern := filepath.Join(pdfDir, baseFilenameWithoutExt+"_noten.pdf")
    files, err := filepath.Glob(pattern)
    if err != nil {
        return fmt.Errorf("failed to glob HTML PDF files: %w", err)
    }
    
    // Verwende bestehende Folder-Patterns oder definiere HTML-spezifische
    folderPatterns := getHTMLFolderPatterns(project)
    
    for _, pdfFile := range files {
        filename := filepath.Base(pdfFile)
        newFilename := fmt.Sprintf("%02d_%s", songIndex, filename)
        
        // Alle HTML-PDFs gehen in das 'noten' Verzeichnis
        // (keine Unterscheidung zwischen A3/A4 wie bei ABC-PDFs)
        targetDir := filepath.Join(outputDir, "druckdateien", "noten")
        
        // Optional: Projekt-spezifische Pattern-Zuordnung
        for pattern, folder := range folderPatterns {
            matched, err := filepath.Match(pattern, filename)
            if err != nil {
                return fmt.Errorf("failed to match HTML pattern: %w", err)
            }
            if matched {
                targetDir = filepath.Join(outputDir, "druckdateien", folder)
                break
            }
        }
        
        err := os.MkdirAll(targetDir, 0755)
        if err != nil {
            return fmt.Errorf("failed to create HTML target directory: %w", err)
        }
        
        targetFile := filepath.Join(targetDir, newFilename)
        err = copyFile(pdfFile, targetFile)
        if err != nil {
            return fmt.Errorf("failed to copy HTML PDF: %w", err)
        }
        
        slog.Info("distributed HTML PDF", "source", pdfFile, "target", targetFile)
    }
    
    return nil
}

func getHTMLFolderPatterns(project *ent.Project) map[string]string {
    // Standard-Patterns für HTML-PDFs - alle gehen in 'noten' Verzeichnis
    // (keine Unterscheidung zwischen A3/A4 wie bei ABC-PDFs)
    patterns := map[string]string{
        "*_noten.pdf": "noten",
    }
    
    // Projekt-spezifische HTML-Patterns laden falls vorhanden
    if configPatterns, ok := project.Config["htmlFolderPatterns"].(map[string]interface{}); ok {
        for pattern, folder := range configPatterns {
            if folderStr, ok := folder.(string); ok {
                patterns[pattern] = folderStr
            }
        }
    }
    
    return patterns
}
```

### 6. Konfiguration und Dependencies

#### go.mod Ergänzung
```go
require (
    github.com/chromedp/chromedp v0.9.3
    // ... bestehende dependencies
)
```

#### Projekt-Konfiguration
```json
{
    "htmlFolderPatterns": {
        "*_noten.pdf": "noten"
    },
    "htmlDOMInjections": [
        {
            "type": "cleanup",
            "action": "remove",
            "selector": "text",
            "pattern": "#vb"
        },
        {
            "type": "pageNumber",
            "position": "bottom-right",
            "style": {
                "background": "grey",
                "padding": "5px",
                "border": "1px solid black",
                "font-weight": "bold"
            }
        },
        {
            "type": "custom",
            "element": "div",
            "id": "customInfo",
            "content": "${SONG_TITLE} - Seite ${SONG_INDEX}",
            "position": "top-left"
        }
    ]
}
```

### 7. Error Handling und Logging

```go
type HTMLConversionError struct {
    Song      string
    HTMLFile  string
    Operation string
    Cause     error
}

func (e *HTMLConversionError) Error() string {
    return fmt.Sprintf("HTML conversion failed for song '%s' (file: %s) during %s: %v", 
        e.Song, e.HTMLFile, e.Operation, e.Cause)
}

// Structured logging für HTML-Operationen
func logHTMLOperation(operation string, song *ent.ProjectSong, details map[string]interface{}) {
    slog.Info("HTML operation", 
        "operation", operation,
        "song", song.Edges.Song.Title,
        "details", details)
}
```

### 8. Testing Strategy

```go
package htmlpdf_test

import (
    "context"
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/bwl21/zupfmanager/internal/htmlpdf"
)

func TestHTMLToPDFConversion(t *testing.T) {
    // Setup test HTML file (simuliert Zupfnoter-Output)
    testHTML := `
    <!DOCTYPE html>
    <html>
    <head><title>Zupfnoter Generated Song</title></head>
    <body>
        <h1>Song Title</h1>
        <div class="music-notation">
            <!-- Zupfnoter-generierte Musik-Notation -->
        </div>
    </body>
    </html>
    `
    
    tempHTML, err := createTempFile(testHTML, "*.html")
    assert.NoError(t, err)
    defer os.Remove(tempHTML)
    
    // Test conversion mit DOM-Injection
    converter := htmlpdf.NewChromeDPConverter(
        htmlpdf.NewTextCleanupInjector("#vb"),
        htmlpdf.NewPageNumberInjector("bottom-right"),
    )
    defer converter.Close()
    
    testSong := createTestSong("Test Song")
    testSong.Edges.Song.Filename = "test_song.abc"
    
    request := &htmlpdf.ConversionRequest{
        HTMLFilePath: tempHTML,
        OutputPath:   "test_song_noten.pdf",  // Entspricht ABC-Dateiname + _noten.pdf
        SongIndex:    5,
        Song:         testSong,
    }
    
    result, err := converter.ConvertToPDF(context.Background(), request)
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Greater(t, result.FileSize, int64(0))
    
    // Validiere dass Original-HTML unverändert ist
    originalContent, err := os.ReadFile(tempHTML)
    assert.NoError(t, err)
    assert.Equal(t, testHTML, string(originalContent))
    
    // Cleanup
    os.Remove(result.OutputPath)
}

func TestDOMInjection(t *testing.T) {
    injector := htmlpdf.NewPageNumberInjector("bottom-right")
    
    request := &htmlpdf.ConversionRequest{
        SongIndex: 42,
        Song:      createTestSong("Test Song"),
        DOMScripts: make([]string, 0),
    }
    
    err := injector.InjectIntoDOM(context.Background(), request)
    assert.NoError(t, err)
    assert.Len(t, request.DOMScripts, 3) // #vb Cleanup + CSS + HTML Element
    
    // Validiere dass #vb Cleanup-Skript vorhanden ist
    foundCleanup := false
    foundPageNumber := false
    for _, script := range request.DOMScripts {
        if strings.Contains(script, "#vb") && strings.Contains(script, "remove") {
            foundCleanup = true
        }
        if strings.Contains(script, "42") {
            foundPageNumber = true
        }
    }
    assert.True(t, foundCleanup, "#vb cleanup script should be present")
    assert.True(t, foundPageNumber, "Page number should be injected into DOM scripts")
}

func TestTextCleanupInjector(t *testing.T) {
    injector := htmlpdf.NewTextCleanupInjector("#vb", "#debug")
    
    request := &htmlpdf.ConversionRequest{
        DOMScripts: make([]string, 0),
    }
    
    err := injector.InjectIntoDOM(context.Background(), request)
    assert.NoError(t, err)
    assert.Len(t, request.DOMScripts, 2) // Zwei Cleanup-Patterns
    
    // Validiere dass beide Patterns im Skript enthalten sind
    script := strings.Join(request.DOMScripts, " ")
    assert.Contains(t, script, "#vb")
    assert.Contains(t, script, "#debug")
    assert.Contains(t, script, "remove()")
}
```

## Implementierungsplan

### Phase 1: Grundstruktur (1-2 Tage)
1. HTML-zu-PDF Interface definieren
2. ChromeDP Basis-Implementation mit DOM-Injection
3. PageNumberInjector implementieren
4. Unit Tests für DOM-Injection

### Phase 2: Integration (2-3 Tage)
1. Integration in buildSong Funktion (parallel zu ABC)
2. Parallele Verarbeitung mit errgroup implementieren
3. PDF-Verteilung für HTML-PDFs
4. Error Handling für DOM-Injection

### Phase 3: Erweiterte Features (1-2 Tage)
1. CustomDOMInjector für flexible Anpassungen
2. Konfigurierbare DOM-Injections über Projekt-Config
3. Projekt-spezifische HTML-Patterns
4. Performance-Optimierungen (Chrome-Pool)

### Phase 4: Testing & Dokumentation (1 Tag)
1. Integration Tests mit echten Zupfnoter-HTML-Dateien
2. Performance Tests (Memory-Usage, Parallelisierung)
3. Dokumentation für DOM-Injection-System
4. Code Review

### Phase 5: Zupfnoter-Kompatibilität (0.5 Tage)
1. Tests mit verschiedenen Zupfnoter-HTML-Outputs
2. CSS-Kompatibilität validieren
3. Print-Styles testen
4. Edge-Cases dokumentieren

## Vorteile dieser Lösung

✅ **Parallelisierung**: ABC und HTML werden gleichzeitig verarbeitet  
✅ **DOM-Injection**: Keine Änderung an Zupfnoter-generierten HTML-Dateien  
✅ **Flexibilität**: Erweiterbare DOM-Injector-Pipeline  
✅ **Isolation**: Original-HTML bleibt vollständig unverändert  
✅ **Integration**: Nahtlose Einbindung in bestehenden Workflow  
✅ **Konfigurierbarkeit**: Projekt-spezifische DOM-Anpassungen möglich  
✅ **Testbarkeit**: Klare Interfaces und Dependency Injection  
✅ **Performance**: ChromeDP ist effizient für PDF-Generierung  
✅ **Zupfnoter-Kompatibilität**: Funktioniert mit allen Zupfnoter-HTML-Outputs  
✅ **Einheitliches Format**: Alle HTML-PDFs in A4, vereinfacht Organisation

## Überlegungen

⚠️ **Chrome Dependency**: ChromeDP benötigt Chrome/Chromium Installation  
⚠️ **Memory Usage**: Chrome kann speicherintensiv sein bei vielen parallelen Konvertierungen  
⚠️ **DOM-Timing**: Kurze Wartezeit nötig damit DOM-Änderungen vor PDF-Generierung wirksam werden  
⚠️ **CSS Print-Styles**: Zupfnoter-CSS könnte mit injiziertem CSS interferieren  
⚠️ **JavaScript**: Falls Zupfnoter JavaScript verwendet, könnte DOM-Injection betroffen sein

## Spezielle Überlegungen für Zupfnoter

✅ **Read-Only HTML**: Perfekt für Zupfnoter-generierte Dateien  
✅ **CSS @media print**: Blattnummer nur beim Drucken/PDF-Export sichtbar  
✅ **Z-Index**: Blattnummer wird über Zupfnoter-Content gelegt  
⚠️ **Zupfnoter-Updates**: Bei Änderungen am Zupfnoter-HTML-Output eventuell Anpassungen nötig

## Beispiel: Konkrete DOM-Injection

### Datei-Namenskonvention:
```
ABC-Datei:    "lied_001.abc"
HTML-Datei:   "lied_001.html"        (von Zupfnoter generiert)
PDF-Output:   "lied_001_noten.pdf"   (HTML-zu-PDF, immer A4)

Verzeichnisstruktur:
druckdateien/
├── klein/          (ABC-PDFs: A, M, O Formate)
├── gross/          (ABC-PDFs: B, X Formate)  
└── noten/          (HTML-PDFs: alle in A4)
```

### Zupfnoter HTML (Original, bleibt unverändert):
```html
<!DOCTYPE html>
<html>
<head>
    <title>Song Title</title>
    <style>/* Zupfnoter CSS */</style>
</head>
<body>
    <div class="music-notation">
        <!-- Zupfnoter-generierte Musik -->
        <text>#vb</text>  <!-- Wird entfernt -->
    </div>
</body>
</html>
```

### DOM-Injection zur Laufzeit:
```javascript
// 1. Entferne <text> Elemente mit '#vb' Inhalt
const textElements = document.querySelectorAll('text');
textElements.forEach(element => {
    if (element.textContent && element.textContent.trim() === '#vb') {
        element.remove();
    }
});

// 2. CSS-Style hinzufügen
const style = document.createElement('style');
style.textContent = `
    @media print {
        #druckParagraph {
            position: fixed;
            bottom: 0;
            right: 0;
            margin: 10px;
            font-weight: bold;
            background: grey;
            padding: 5px;
            border: 1px solid black;
            z-index: 1000;
        }
    }
`;
document.head.appendChild(style);

// 3. HTML-Element hinzufügen
const paragraph = document.createElement('p');
paragraph.id = 'druckParagraph';
paragraph.textContent = '05'; // Blattnummer
document.body.insertBefore(paragraph, document.body.firstChild);
```

### Resultat im PDF:
- **Dateiname**: `lied_001_noten.pdf` (basierend auf ABC-Dateiname)
- **Format**: Immer A4 (keine A3/A4 Unterscheidung wie bei ABC-PDFs)
- **Verzeichnis**: Alle HTML-PDFs gehen in `druckdateien/noten/`
- **#vb Text-Elemente werden entfernt** (nicht mehr im PDF sichtbar)
- Original Zupfnoter-Musik-Notation bleibt unverändert
- Blattnummer "05" erscheint unten rechts beim PDF-Export
- Blattnummer ist nur im Print-Modus sichtbar (@media print)
- Original HTML-Datei bleibt vollständig unverändert

---

*Diese Spezifikation bietet eine vollständige Lösung für HTML-zu-PDF Konvertierung parallel zur bestehenden ABC-Verarbeitung, speziell angepasst für Zupfnoter-generierte HTML-Dateien.*