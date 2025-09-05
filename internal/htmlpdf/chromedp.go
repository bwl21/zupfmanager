package htmlpdf

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// ChromeDPConverter implements HTMLToPDFConverter using ChromeDP
type ChromeDPConverter struct {
	allocCtx   context.Context
	cancelFunc context.CancelFunc
	injectors  []DOMInjector
}

// NewChromeDPConverter creates a new ChromeDP-based HTML to PDF converter
func NewChromeDPConverter(injectors ...DOMInjector) *ChromeDPConverter {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.NoSandbox,
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-web-security", true), // For local files
	)

	return &ChromeDPConverter{
		allocCtx:   allocCtx,
		cancelFunc: cancel,
		injectors:  injectors,
	}
}

// ConvertToPDF converts HTML to PDF using ChromeDP with DOM injection
func (c *ChromeDPConverter) ConvertToPDF(ctx context.Context, request *ConversionRequest) (*ConversionResult, error) {
	start := time.Now()

	// 1. Validate that HTML file exists (Zupfnoter-generated)
	if _, err := os.Stat(request.HTMLFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("HTML file does not exist: %s", request.HTMLFilePath)
	}

	// 2. Prepare DOM injectors
	request.DOMScripts = make([]string, 0)
	for _, injector := range c.injectors {
		err := injector.InjectIntoDOM(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("DOM injector %s failed: %w", injector.Name(), err)
		}
	}

	// 3. Generate PDF with DOM manipulation
	var pdfBuffer []byte
	taskCtx, cancel := chromedp.NewContext(c.allocCtx)
	defer cancel()

	// Create ChromeDP actions
	actions := []chromedp.Action{
		chromedp.Navigate("file://" + request.HTMLFilePath),
		chromedp.WaitReady("body"),
	}

	// Add DOM manipulation scripts
	for _, script := range request.DOMScripts {
		actions = append(actions, chromedp.Evaluate(script, nil))
	}

	// Wait briefly for DOM changes to take effect
	actions = append(actions, chromedp.Sleep(500*time.Millisecond))

	// PDF generation (fixed A4 format for HTML PDFs)
	actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		pdfBuffer, _, err = page.PrintToPDF().
			WithPrintBackground(true).
			WithPaperWidth(8.27).  // A4 width in inches (fixed)
			WithPaperHeight(11.7). // A4 height in inches (fixed)
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

	// 4. Write PDF file
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

// ValidateHTML validates that the HTML file exists and is readable
func (c *ChromeDPConverter) ValidateHTML(htmlPath string) error {
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		return fmt.Errorf("HTML file does not exist: %s", htmlPath)
	}
	return nil
}

// Close cleans up the ChromeDP allocator
func (c *ChromeDPConverter) Close() error {
	if c.cancelFunc != nil {
		c.cancelFunc()
	}
	return nil
}