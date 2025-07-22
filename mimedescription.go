package mimedescription

// Get returns the human-friendly description for a given MIME type.
// The second return value is false if the MIME type is not found.
func Get(mimeType string) (string, bool) {
	description, ok := mimeData[mimeType]
	return description, ok
}
