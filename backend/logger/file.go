package logger

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/meh-hackathon/meh/apperror"
	"github.com/meh-hackathon/meh/constants"
)

var (
	ErrRotation        = apperror.Define("logger:rotation", "Error during log rotation")
	ErrOpenLogFile     = apperror.Define("logger:open_file", "Error opening log file")
	ErrSetupFileLogger = apperror.Define("logger:setup_file", "Error setting up file logger")
	ErrCloseLogFile    = apperror.Define("logger:close_file", "Error closing log file")
)

var _ Logger = (*FileLogger)(nil)

type FileLogger struct {
	level    slog.Level
	filename string

	rotationConfig RotationConfig

	file     *os.File
	fileSize int64
	mu       sync.Mutex
}

type RotationConfig struct {
	MaxSize       int64
	MaxFiles      int
	Compress      bool
	RotateOnStart bool
}

type FileLoggerOption func(*FileLogger)

func WithRotation(cfg RotationConfig) FileLoggerOption {
	return func(logger *FileLogger) { logger.rotationConfig = cfg }
}

func WithLevel(level slog.Level) FileLoggerOption {
	return func(logger *FileLogger) { logger.level = level }
}

func NewFileLogger(filename string, options ...FileLoggerOption) (*FileLogger, error) {
	logger := &FileLogger{
		filename: filename,
		level:    slog.LevelInfo,
		rotationConfig: RotationConfig{
			MaxSize:       10 * constants.MB,
			MaxFiles:      5,
			Compress:      true,
			RotateOnStart: false,
		},
	}

	for _, option := range options {
		option(logger)
	}

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, ErrSetupFileLogger.WithMessage("Failed to create log directory").WithOrigin().WithCause(err)
	}

	if logger.rotationConfig.RotateOnStart {
		_, err := os.Stat(filename)
		if err != nil && !os.IsNotExist(err) {
			return nil, ErrSetupFileLogger.WithMessage("Failed to check log file existence").WithOrigin().WithCause(err)
		}
		if err == nil {
			if err := logger.rotate(); err != nil {
				return nil, ErrSetupFileLogger.WithMessage("Failed to rotate existing log file").WithOrigin().WithCause(err)
			}
		}
	}

	if err := logger.openFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (fl *FileLogger) log(level slog.Level, msg string, args ...any) {
	if level < fl.level {
		return
	}
	fl.write(level, time.Now(), msg, args...)
}

func (fl *FileLogger) Debug(msg string, args ...any) { fl.log(slog.LevelDebug, msg, args...) }
func (fl *FileLogger) Info(msg string, args ...any)  { fl.log(slog.LevelInfo, msg, args...) }
func (fl *FileLogger) Warn(msg string, args ...any)  { fl.log(slog.LevelWarn, msg, args...) }
func (fl *FileLogger) Error(msg string, args ...any) { fl.log(slog.LevelError, msg, args...) }

func (fl *FileLogger) openFile() error {
	file, err := os.OpenFile(fl.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return ErrOpenLogFile.WithMessage("Failed to open log file").WithOrigin().WithCause(err)
	}

	stat, err := file.Stat()
	if err != nil {
		return ErrOpenLogFile.WithMessage("Failed to get log file stats").WithOrigin().WithCause(err)
	}

	fl.file = file
	fl.fileSize = stat.Size()
	return nil
}

func (fl *FileLogger) write(level slog.Level, timestamp time.Time, msg string, args ...any) {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	if fl.needsRotation() {
		if err := fl.rotate(); err != nil {
			// Handle rotation error - could log to stderr or return
			slog.Error("Failed to rotate log file", "error", err)
			return
		}
	}

	logEntry := fl.formatAsJSONL(level, timestamp, msg, args...)
	addedSize, err := fl.file.WriteString(logEntry)
	if err != nil {
		slog.Error("Failed to write log entry", "error", err, "message", msg, "args", args)
		return
	}
	fl.fileSize += int64(addedSize)
}

func (fl *FileLogger) needsRotation() bool {
	if fl.rotationConfig.MaxFiles <= 0 {
		return false
	}
	if fl.fileSize < fl.rotationConfig.MaxSize {
		return false
	}

	return true
}

func (fl *FileLogger) rotate() error {
	if fl.file != nil {
		if err := fl.file.Close(); err != nil {
			return ErrRotation.WithMessage("Failed to close log file").WithOrigin().WithCause(ErrCloseLogFile.WithCause(err))
		}
	}

	timestamp := time.Now().Format("2006-01-02-15-04-05")
	rotatedName := fmt.Sprintf("%s.%s", fl.filename, timestamp)

	if err := os.Rename(fl.filename, rotatedName); err != nil {
		return ErrRotation.WithMessage("Failed to rename log file").WithOrigin().WithCause(err)
	}

	if fl.rotationConfig.Compress {
		go fl.compressFile(rotatedName)
	}

	go fl.cleanupOldFiles()

	return fl.openFile()
}

func (fl *FileLogger) compressFile(filename string) {
	// Create the compressed filename
	compressedName := filename + ".gz"

	// Open the source file
	srcFile, err := os.Open(filename)
	if err != nil {
		slog.Error("Failed to open file for compression", "file", filename, "error", err)
		return
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(compressedName)
	if err != nil {
		slog.Error("Failed to create compressed file", "file", compressedName, "error", err)
		return
	}
	defer dstFile.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(dstFile)
	defer gzipWriter.Close()

	// Copy the content through gzip writer
	_, err = io.Copy(gzipWriter, srcFile)
	if err != nil {
		slog.Error("Failed to compress file content", "file", filename, "error", err)
		os.Remove(compressedName)
		return
	}

	if err := gzipWriter.Close(); err != nil {
		slog.Error("Failed to close gzip writer", "file", filename, "error", err)
		os.Remove(compressedName)
		return
	}

	if err := dstFile.Close(); err != nil {
		slog.Error("Failed to close compressed file", "file", compressedName, "error", err)
		os.Remove(compressedName)
		return
	}

	if err := os.Remove(filename); err != nil {
		slog.Error("Failed to remove original file after compression", "file", filename, "error", err)
		return
	}
}

// cleanupOldFiles removes old rotated log files, keeping only MaxFiles most recent ones
func (fl *FileLogger) cleanupOldFiles() {
	if fl.rotationConfig.MaxFiles <= 0 {
		return
	}

	dir := filepath.Dir(fl.filename)
	baseName := filepath.Base(fl.filename)

	entries, err := os.ReadDir(dir)
	if err != nil {
		slog.Error("Failed to read log directory for cleanup", "dir", dir, "error", err)
		return
	}

	// Create regex pattern to match rotated log files
	// Pattern: {basename}.{timestamp} or {basename}.{timestamp}.gz
	escapedBaseName := regexp.QuoteMeta(baseName)
	pattern := fmt.Sprintf(`^%s\.\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}(\.gz)?$`, escapedBaseName)
	regex, err := regexp.Compile(pattern)
	if err != nil {
		slog.Error("Failed to compile regex for log file cleanup", "pattern", pattern, "error", err)
		return
	}

	// Find all matching rotated log files
	var rotatedFiles []rotatedLogFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if regex.MatchString(filename) {
			info, err := entry.Info()
			if err != nil {
				slog.Error("Failed to get file info during cleanup", "file", filename, "error", err)
				continue
			}

			// Extract timestamp from filename for sorting
			timestamp := fl.extractTimestampFromFilename(filename, baseName)
			rotatedFiles = append(rotatedFiles, rotatedLogFile{
				name:      filename,
				timestamp: timestamp,
				modTime:   info.ModTime(),
				fullPath:  filepath.Join(dir, filename),
			})
		}
	}

	// Sort by timestamp (newest first), fallback to modification time
	sort.Slice(rotatedFiles, func(i, j int) bool {
		if rotatedFiles[i].timestamp != rotatedFiles[j].timestamp {
			return rotatedFiles[i].timestamp > rotatedFiles[j].timestamp
		}
		return rotatedFiles[i].modTime.After(rotatedFiles[j].modTime)
	})

	// Remove files beyond MaxFiles limit
	if len(rotatedFiles) > fl.rotationConfig.MaxFiles {
		filesToRemove := rotatedFiles[fl.rotationConfig.MaxFiles:]
		for _, file := range filesToRemove {
			if err := os.Remove(file.fullPath); err != nil {
				slog.Error("Failed to remove old log file", "file", file.fullPath, "error", err)
			} else {
				slog.Debug("Removed old log file", "file", file.name)
			}
		}
	}
}

type rotatedLogFile struct {
	name      string
	timestamp string
	modTime   time.Time
	fullPath  string
}

func (fl *FileLogger) extractTimestampFromFilename(filename, baseName string) string {
	remainder := strings.TrimPrefix(filename, baseName+".")
	remainder = strings.TrimSuffix(remainder, ".gz")
	return remainder
}

func (fl *FileLogger) Close() error {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	if fl.file != nil {
		if err := fl.file.Close(); err != nil {
			return ErrCloseLogFile.WithOrigin().WithCause(err)
		}
	}
	return nil
}

func (fl *FileLogger) formatAsJSONL(level slog.Level, timestamp time.Time, msg string, args ...any) string {
	var data map[string]any
	if len(args) > 0 {
		data = make(map[string]any)

		for i := 0; i < len(args)-1; i += 2 {
			key := fmt.Sprintf("%v", args[i])
			data[key] = args[i+1]
		}

		if len(args)%2 == 1 {
			data["_extra"] = args[len(args)-1]
		}
	}

	entry := fmt.Sprintf(`{"timestamp":"%s", "level":"%s", "message":"%s"`, timestamp.Format(time.RFC3339Nano), level.String(), msg)
	if data != nil {
		dataStr, err := json.Marshal(data)
		if err != nil {
			slog.Error("Failed to marshal additional data", "error", err)
			dataStr = []byte(`{"error":"failed to marshal additional data"}`)
		}
		entry += fmt.Sprintf(`, "data":%s`, dataStr)
	}
	entry += "}\n"

	return entry
}

func (fl *FileLogger) SetLevel(level slog.Level) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.level = level
}
