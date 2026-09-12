package runtime

import "strings"

// MaxTextChars is the recorded maximum number of characters in one string.
const MaxTextChars = 255

// LimitText keeps the first 255 characters of a reference string value.
// Decoded Windows-1252 characters occupy one position each; Unicode extensions
// use the same character limit without splitting UTF-8 encodings.
func LimitText(text string) string {
	count := 0
	for offset := range text {
		if count == MaxTextChars {
			return strings.Clone(text[:offset])
		}
		count++
	}
	return text
}
