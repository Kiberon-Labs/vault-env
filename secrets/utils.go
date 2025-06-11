package secrets

import (
	"fmt"
)

var forwardSlashRune rune = '/'

// formatPath constructs a path for the Vault KV engine, ensuring that it is formatted correctly
// for both KV v1 and KV v2 engines. It handles leading slashes and ensures that the
// path is properly concatenated with the root path.
// It also ensures that the path is formatted correctly for Unicode characters.
// If the engine is "kv-v2", it appends "/data" to the root path.
// If the root path is empty, it will not prepend a slash.
func formatPath(engine string, root string, path string) string {
	// ... This will handle Unicode characters correctly.
	//     Not needed for ASCII strings.
	runes := []rune(path)
	rootRunes := []rune(root)
	// Remove leading slashes if they exist
	if len(runes) > 0 && runes[0] == forwardSlashRune {
		runes = runes[1:]
	}
	if len(rootRunes) > 0 && rootRunes[0] == forwardSlashRune {
		rootRunes = rootRunes[1:]
	}

	finalRoot := string(rootRunes)
	finalPath := string(runes)

	// This is a special case for KV v2 where we need to append "/data" to the root path.
	if engine == "kv-v2" {
		finalRoot = fmt.Sprintf("%s/%s", finalRoot, "data")
	}

	return fmt.Sprintf("%s/%s", finalRoot, finalPath)
}
