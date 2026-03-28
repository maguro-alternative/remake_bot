package htmlutil

import "html"

// EscapeString escapes special characters in a string for safe HTML embedding.
// This should be used for all user-controlled or external data before embedding in HTML.
func EscapeString(s string) string {
	return html.EscapeString(s)
}
