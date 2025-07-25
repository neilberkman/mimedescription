package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

const (
	// URL for the authoritative freedesktop.org shared-mime-info database.
	dbURL = "https://gitlab.freedesktop.org/xdg/shared-mime-info/-/raw/master/data/freedesktop.org.xml.in"
	// Path to the output file for the generated Go code.
	outputPath = "data.go"
	// Template for the generated Go file.
	fileTemplate = `// Code generated by go run cmd/generator/main.go; DO NOT EDIT.
// This file was generated on {{ .Timestamp }}
package mimedescription

// mimeData is a map of MIME types to their human-friendly descriptions.
// It is sourced from the freedesktop.org shared-mime-info database.
var mimeData = map[string]string{
{{- range $mime, $desc := .MimeMap }}
	"{{ $mime }}": "{{ $desc }}",
{{- end }}
}
`
)

// MimeType represents a <mime-type> element in the XML.
type MimeType struct {
	Type    string `xml:"type,attr"`
	Comment string `xml:"comment"`
}

func main() {
	log.Println("Starting MIME description generator...")

	// 1. Fetch the XML data
	log.Printf("Fetching data from %s...", dbURL)
	resp, err := http.Get(dbURL)
	if err != nil {
		log.Fatalf("Failed to fetch MIME database: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch MIME database: received status code %d", resp.StatusCode)
	}

	xmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// 2. Parse the XML data
	log.Println("Parsing XML data...")
	descriptions, err := parseMimeXML(bytes.NewReader(xmlData))
	if err != nil {
		log.Fatalf("Failed to parse MIME XML: %v", err)
	}
	log.Printf("Successfully parsed %d MIME descriptions.", len(descriptions))

	// 3. Generate the Go source file from the template
	log.Printf("Generating Go source file at %s...", outputPath)
	if err := generateGoFile(descriptions); err != nil {
		log.Fatalf("Failed to generate Go file: %v", err)
	}

	log.Println("Generator finished successfully.")
}

// parseMimeXML reads the XML data and extracts a map of MIME types to comments.
func parseMimeXML(reader io.Reader) (map[string]string, error) {
	decoder := xml.NewDecoder(reader)
	descriptions := make(map[string]string)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error decoding token: %w", err)
		}

		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "mime-type" {
			var mt MimeType
			if err := decoder.DecodeElement(&mt, &se); err != nil {
				return nil, fmt.Errorf("error decoding mime-type element: %w", err)
			}
			if mt.Comment != "" {
				// The descriptions can have quotes, so we need to escape them.
				escapedComment := escapeString(mt.Comment)
				descriptions[mt.Type] = escapedComment
			}
		}
	}
	return descriptions, nil
}

// generateGoFile takes the map of descriptions and writes the data.go file.
func generateGoFile(descriptions map[string]string) error {
	// Prepare data for the template
	data := struct {
		Timestamp string
		MimeMap   map[string]string
	}{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		MimeMap:   descriptions,
	}

	// Create a new template and parse the file template
	tmpl, err := template.New("data").Parse(fileTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template and write to a buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Format the generated code using go/format
	formattedSource, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format generated code: %w", err)
	}

	// Write the formatted code to the output file
	err = ioutil.WriteFile(outputPath, formattedSource, 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// escapeString escapes backslashes and double quotes for inclusion in a Go string literal.
func escapeString(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			result.WriteString(`\\`)
		case '"':
			result.WriteString(`\"`)
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}
