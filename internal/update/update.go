// Package update provides update checking functionality.
package update

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// GitHub API endpoint for latest release
	githubReleaseURL = "https://api.github.com/repos/gurselcakar/arithmego/releases/latest"
	// GitHub download URL pattern for release assets
	githubDownloadURL = "https://github.com/gurselcakar/arithmego/releases/download/%s/arithmego_%s_%s.%s"
	// Timeout for HTTP requests (check)
	httpTimeout = 10 * time.Second
	// Timeout for download requests (longer for large files)
	downloadTimeout = 60 * time.Second
)

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

// DownloadAndApply downloads the release archive for the given version and
// replaces the running binary. The caller should prompt the user to restart.
func DownloadAndApply(version string) error {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	url := fmt.Sprintf(githubDownloadURL, version, goos, goarch, "tar.gz")

	archivePath, err := downloadArchive(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(archivePath)

	binaryPath, err := extractFromTarGz(archivePath)
	if err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}
	defer os.Remove(binaryPath)

	if err := replaceBinary(binaryPath); err != nil {
		return fmt.Errorf("replace failed: %w", err)
	}

	return nil
}

// downloadArchive downloads the release archive to a temporary file.
func downloadArchive(url string) (string, error) {
	client := &http.Client{Timeout: downloadTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "arithmego-update-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write archive: %w", err)
	}

	return tmpFile.Name(), nil
}

// extractFromTarGz extracts the arithmego binary from a tar.gz archive.
func extractFromTarGz(archivePath string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", fmt.Errorf("failed to open gzip: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read tar: %w", err)
		}

		name := filepath.Base(header.Name)
		if name == "arithmego" && header.Typeflag == tar.TypeReg {
			return writeToTemp(tr)
		}
	}

	return "", fmt.Errorf("arithmego binary not found in archive")
}

// writeToTemp writes from a reader to a temporary file with executable permissions.
func writeToTemp(r io.Reader) (string, error) {
	tmpFile, err := os.CreateTemp("", "arithmego-bin-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, r); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	if err := os.Chmod(tmpFile.Name(), 0o755); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

// replaceBinary replaces the current running binary with the new one.
func replaceBinary(newBinaryPath string) error {
	currentPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Resolve symlinks to get the actual binary path
	currentPath, err = filepath.EvalSymlinks(currentPath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	// Try direct rename first (fastest, works when on same filesystem)
	if err := os.Rename(newBinaryPath, currentPath); err != nil {
		// Cross-device rename: copy to same directory then rename
		dir := filepath.Dir(currentPath)
		tmpDest, err2 := os.CreateTemp(dir, "arithmego-new-*")
		if err2 != nil {
			return fmt.Errorf("rename failed (%v) and temp create failed: %w", err, err2)
		}
		tmpDestPath := tmpDest.Name()

		src, err2 := os.Open(newBinaryPath)
		if err2 != nil {
			os.Remove(tmpDestPath)
			return fmt.Errorf("rename failed (%v) and open failed: %w", err, err2)
		}

		if _, err2 = io.Copy(tmpDest, src); err2 != nil {
			src.Close()
			tmpDest.Close()
			os.Remove(tmpDestPath)
			return fmt.Errorf("rename failed (%v) and copy failed: %w", err, err2)
		}
		src.Close()
		tmpDest.Close()

		if err2 = os.Chmod(tmpDestPath, 0o755); err2 != nil {
			os.Remove(tmpDestPath)
			return fmt.Errorf("chmod failed: %w", err2)
		}

		if err2 = os.Rename(tmpDestPath, currentPath); err2 != nil {
			os.Remove(tmpDestPath)
			return fmt.Errorf("final rename failed: %w", err2)
		}
	}

	return nil
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
