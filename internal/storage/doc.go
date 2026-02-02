// Package storage handles persistence of user data to the local filesystem.
//
// This package manages two types of data:
//
//   - Configuration ([Config]): User preferences, defaults, and quick play state.
//     Stored in config.json. Non-critical data that falls back to defaults on error.
//
//   - Statistics ([Statistics]): Game session history with detailed question records.
//     Stored in statistics.json. Critical data that returns errors on corruption.
//
// All files are stored in the user's config directory under "arithmego".
// Use [ConfigDir] to get the directory path, or [ConfigPath] and [StatisticsPath]
// for specific file paths.
//
// Both [SaveConfig] and [Save] use atomic writes (write to temp file, then rename)
// to prevent data corruption on crashes or power loss.
package storage
