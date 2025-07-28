package extractid

import "strings"

// extractUserIDFromPath extracts user ID from URL paths like example.domain/followers/:id
func ExtractUserIDFromPath(path, endpoint string) string {
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, endpoint+"/")
	return parts[len(parts)-1]
}
