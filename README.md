# Go MIME Description Library

A Go library that provides human-friendly descriptions for MIME types, automatically kept up-to-date with the freedesktop.org shared-mime-info database.

## Features

- 🚀 **Fast**: Pre-built map with no runtime parsing
- 🔄 **Always up-to-date**: Automatically synced weekly via GitHub Actions
- 📦 **Zero dependencies**: Pure Go with no external dependencies
- 🎯 **Simple API**: Just one function to get descriptions

## Installation

```bash
go get github.com/neilberkman/mimedescription
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/neilberkman/mimedescription"
)

func main() {
    // Get description for a MIME type
    description, found := mimedescription.Get("application/pdf")
    if found {
        fmt.Println(description) // Output: "PDF document"
    } else {
        fmt.Println("MIME type not found")
    }

    // Example with other MIME types
    examples := []string{
        "image/jpeg",
        "text/html",
        "application/json",
        "video/mp4",
        "audio/mpeg",
        "unknown/type",
    }

    for _, mimeType := range examples {
        if desc, ok := mimedescription.Get(mimeType); ok {
            fmt.Printf("%s -> %s\n", mimeType, desc)
        } else {
            fmt.Printf("%s -> (unknown)\n", mimeType)
        }
    }
}
```

Output:
```
PDF document
image/jpeg -> JPEG image
text/html -> HTML document  
application/json -> JSON document
video/mp4 -> MPEG-4 video
audio/mpeg -> MP3 audio
unknown/type -> (unknown)
```

## API Reference

### `Get(mimeType string) (string, bool)`

Returns the human-friendly description for the given MIME type.

**Parameters:**
- `mimeType`: The MIME type string (e.g., "application/pdf")

**Returns:**
- `string`: Human-friendly description of the MIME type
- `bool`: `true` if the MIME type was found, `false` otherwise

## Data Source

This library uses the **freedesktop.org shared-mime-info database**, which contains **955 MIME types** covering the most commonly used file formats. This includes:

- Document formats (PDF, Word, etc.)
- Image formats (JPEG, PNG, WebP, etc.)
- Audio formats (MP3, FLAC, Ogg, etc.)
- Video formats (MP4, WebM, AVI, etc.)
- Archive formats (ZIP, TAR, 7z, etc.)
- Programming languages and scripts
- And many more...

**Note:** This is not a complete list of all MIME types that exist. The IANA registry contains thousands of registered MIME types, and applications can define custom ones. This library focuses on the most commonly encountered types that users are likely to interact with.

## Automatic Updates

The MIME type data is automatically updated weekly via GitHub Actions. The workflow:

1. Fetches the latest XML data from freedesktop.org
2. Regenerates the Go source file
3. Commits changes if the data has been updated

This ensures your application always has the latest MIME type definitions without requiring manual updates.

## Contributing

Contributions are welcome! The main areas for contribution:

1. **Bug fixes**: If you find issues with the generator or library
2. **Feature requests**: Additional functionality that would be useful
3. **Data issues**: Report if specific MIME types are missing or incorrectly described

Note that the MIME type data itself comes from freedesktop.org, so data changes should be contributed upstream to the shared-mime-info project.

## License

This project is open source. The MIME type data is sourced from the freedesktop.org shared-mime-info database.

## Similar Projects

If you need more comprehensive MIME type support, consider:
- Libraries that include the full IANA registry
- File type detection libraries that analyze file content
- Custom MIME type databases for specialized applications