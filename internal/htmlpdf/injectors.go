package htmlpdf

import (
	"context"
	"fmt"
	"strings"
)

// PageNumberInjector injects page numbers into HTML documents
type PageNumberInjector struct {
	cssStyle string
	position string // "top-right", "bottom-right", etc.
}

// NewPageNumberInjector creates a new page number injector
func NewPageNumberInjector(position string) *PageNumberInjector {
	cssStyle := `
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
    `

	return &PageNumberInjector{
		cssStyle: cssStyle,
		position: position,
	}
}

// InjectIntoDOM injects page number CSS and HTML element into the DOM
func (inj *PageNumberInjector) InjectIntoDOM(ctx context.Context, request *ConversionRequest) error {
	pageNumber := fmt.Sprintf("%02d", request.SongIndex)

	// 1. Remove <text> elements with '#vb' content
	removeVbScript := `
        const textElements = document.querySelectorAll('text');
        textElements.forEach(element => {
            if (element.textContent && element.textContent.trim() === '#vb') {
                element.remove();
            }
        });
    `

	// 2. Add CSS style to <head>
	styleScript := fmt.Sprintf(`
        const style = document.createElement('style');
        style.textContent = %s;
        document.head.appendChild(style);
    `, "`"+inj.cssStyle+"`")

	// 3. Add HTML element at the beginning of <body>
	elementScript := fmt.Sprintf(`
        const paragraph = document.createElement('p');
        paragraph.id = 'druckParagraph';
        paragraph.textContent = '%s';
        document.body.insertBefore(paragraph, document.body.firstChild);
    `, pageNumber)

	// All scripts will be executed by ChromeDP
	request.DOMScripts = append(request.DOMScripts, removeVbScript, styleScript, elementScript)

	return nil
}

// Name returns the name of this injector
func (inj *PageNumberInjector) Name() string {
	return "PageNumberInjector"
}

// TextCleanupInjector removes specific text elements from HTML documents
type TextCleanupInjector struct {
	removePatterns []string
}

// NewTextCleanupInjector creates a new text cleanup injector
func NewTextCleanupInjector(patterns ...string) *TextCleanupInjector {
	if len(patterns) == 0 {
		patterns = []string{"#vb"} // Default pattern
	}
	return &TextCleanupInjector{
		removePatterns: patterns,
	}
}

// InjectIntoDOM removes text elements matching the specified patterns
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

// Name returns the name of this injector
func (inj *TextCleanupInjector) Name() string {
	return "TextCleanupInjector"
}

// CustomDOMInjector provides flexible DOM manipulation capabilities
type CustomDOMInjector struct {
	cssStyles    []string
	htmlElements []HTMLElement
	cleanupRules []CleanupRule
}

// NewCustomDOMInjector creates a new custom DOM injector
func NewCustomDOMInjector() *CustomDOMInjector {
	return &CustomDOMInjector{
		cssStyles:    make([]string, 0),
		htmlElements: make([]HTMLElement, 0),
		cleanupRules: make([]CleanupRule, 0),
	}
}

// AddCSS adds a CSS style to be injected
func (inj *CustomDOMInjector) AddCSS(css string) *CustomDOMInjector {
	inj.cssStyles = append(inj.cssStyles, css)
	return inj
}

// AddElement adds an HTML element to be injected
func (inj *CustomDOMInjector) AddElement(element HTMLElement) *CustomDOMInjector {
	inj.htmlElements = append(inj.htmlElements, element)
	return inj
}

// AddCleanupRule adds a cleanup rule to be applied
func (inj *CustomDOMInjector) AddCleanupRule(rule CleanupRule) *CustomDOMInjector {
	inj.cleanupRules = append(inj.cleanupRules, rule)
	return inj
}

// InjectIntoDOM applies all configured DOM manipulations
func (inj *CustomDOMInjector) InjectIntoDOM(ctx context.Context, request *ConversionRequest) error {
	// 1. Apply cleanup rules (first)
	for _, rule := range inj.cleanupRules {
		cleanupScript := inj.generateCleanupScript(rule)
		request.DOMScripts = append(request.DOMScripts, cleanupScript)
	}

	// 2. Add CSS styles
	for _, css := range inj.cssStyles {
		styleScript := fmt.Sprintf(`
            const style = document.createElement('style');
            style.textContent = %s;
            document.head.appendChild(style);
        `, "`"+css+"`")
		request.DOMScripts = append(request.DOMScripts, styleScript)
	}

	// 3. Add HTML elements
	for _, element := range inj.htmlElements {
		elementScript := inj.generateElementScript(element, request)
		request.DOMScripts = append(request.DOMScripts, elementScript)
	}

	return nil
}

// generateCleanupScript generates JavaScript for cleanup rules
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

// generateElementScript generates JavaScript for HTML element injection
func (inj *CustomDOMInjector) generateElementScript(element HTMLElement, request *ConversionRequest) string {
	// Replace placeholders
	content := element.Content
	content = strings.ReplaceAll(content, "${SONG_INDEX}", fmt.Sprintf("%02d", request.SongIndex))
	if request.Song != nil && request.Song.Edges.Song != nil {
		content = strings.ReplaceAll(content, "${SONG_TITLE}", request.Song.Edges.Song.Title)
	}

	script := fmt.Sprintf(`
        const element = document.createElement('%s');
        element.id = '%s';
        element.className = '%s';
        element.textContent = '%s';
    `, element.Tag, element.ID, element.Class, content)

	// Add attributes
	for key, value := range element.Attributes {
		script += fmt.Sprintf(`element.setAttribute('%s', '%s');`, key, value)
	}

	// Determine position
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

// Name returns the name of this injector
func (inj *CustomDOMInjector) Name() string {
	return "CustomDOMInjector"
}