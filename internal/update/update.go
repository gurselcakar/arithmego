// Package update provides update checking functionality.
package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// GitHub API endpoint for latest release
	githubReleaseURL = "https://api.github.com/repos/gurselcakar/arithmego/releases/latest"
	// Timeout for HTTP requests
	httpTimeout = 10 * time.Second
)

// TODO: Implement actual auto-update functionality in Phase 12 (Distribution).
// This would involve:
// - Downloading the appropriate binary for the user's OS/arch
// - Verifying checksums for security
// - Replacing the running binary (platform-specific)
// - Handling permissions and restart
// For now, we only CHECK for updates and notify the user to run 'arithmego update'.

// Info contains information about an available update.
type Info struct {
	CurrentVersion  string
	LatestVersion   string
	ReleaseURL      string
	UpdateAvailable bool
}

// Check checks GitHub for a newer release.
func Check(currentVersion string) (*Info, error) {
	info := &Info{
		CurrentVersion: currentVersion,
	}

	// Dev builds can't be compared
	if currentVersion == "dev" {
		info.LatestVersion = "unknown (dev build)"
		info.UpdateAvailable = false
		return info, nil
	}

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(githubReleaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to contact GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// No releases yet
		info.LatestVersion = "none (no releases yet)"
		info.UpdateAvailable = false
		return info, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	if release.TagName == "" {
		return nil, fmt.Errorf("GitHub API returned empty tag_name")
	}

	info.LatestVersion = release.TagName
	info.ReleaseURL = release.HTMLURL
	info.UpdateAvailable = IsNewerVersion(release.TagName, currentVersion)

	return info, nil
}

// IsNewerVersion compares two semver-style version strings.
// Returns true if latest is newer than current.
func IsNewerVersion(latest, current string) bool {
	// Strip 'v' prefix if present
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	// Compare each part numerically
	for i := 0; i < len(latestParts) && i < len(currentParts); i++ {
		var latestNum, currentNum int
		latestParsed, _ := fmt.Sscanf(latestParts[i], "%d", &latestNum)
		currentParsed, _ := fmt.Sscanf(currentParts[i], "%d", &currentNum)

		// If either part failed to parse as a number, fall back to string comparison
		if latestParsed == 0 || currentParsed == 0 {
			if latestParts[i] > currentParts[i] {
				return true
			}
			if latestParts[i] < currentParts[i] {
				return false
			}
			continue
		}

		if latestNum > currentNum {
			return true
		}
		if latestNum < currentNum {
			return false
		}
	}

	// If all compared parts are equal, longer version is newer (e.g., 1.0.1 > 1.0)
	return len(latestParts) > len(currentParts)
}
