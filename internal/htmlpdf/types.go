package htmlpdf

import (
	"context"
	"time"

	"github.com/bwl21/zupfmanager/internal/ent"
)

// HTMLToPDFConverter defines the interface for converting HTML to PDF
type HTMLToPDFConverter interface {
	ConvertToPDF(ctx context.Context, request *ConversionRequest) (*ConversionResult, error)
	ValidateHTML(htmlPath string) error
	Close() error
}

// ConversionRequest contains all parameters needed for HTML to PDF conversion
type ConversionRequest struct {
	HTMLFilePath string            // Path to the HTML source file (Zupfnoter-generated)
	OutputPath   string            // Target PDF path
	SongIndex    int               // Page number for DOM injection
	Song         *ent.ProjectSong  // Song information
	DOMInjectors []DOMInjector     // List of DOM injectors
	DOMScripts   []string          // JavaScript code for DOM manipulation
}

// ConversionResult contains the results of HTML to PDF conversion
type ConversionResult struct {
	OutputPath string        // Path to generated PDF
	PageCount  int           // Number of pages
	FileSize   int64         // File size in bytes
	Duration   time.Duration // Conversion time
	Warnings   []string      // Warnings during conversion
}

// DOMInjector defines the interface for DOM manipulation during conversion
type DOMInjector interface {
	InjectIntoDOM(ctx context.Context, request *ConversionRequest) error
	Name() string
}

// HTMLElement represents an HTML element to be injected
type HTMLElement struct {
	Tag        string
	ID         string
	Class      string
	Content    string
	Attributes map[string]string
	Position   InsertPosition
}

// InsertPosition defines where to insert HTML elements
type InsertPosition int

const (
	BodyStart InsertPosition = iota
	BodyEnd
	HeadEnd
	BeforeElement
	AfterElement
)

// CleanupRule defines rules for cleaning up HTML content
type CleanupRule struct {
	Selector string // CSS selector for elements
	Action   string // "remove", "hide", "modify"
	Pattern  string // Text pattern for matching
	Value    string // New value for "modify" action
}