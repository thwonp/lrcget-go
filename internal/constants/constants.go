package constants

import "time"

// Application constants
const (
	AppName        = "LRCGET"
	AppVersion     = "1.0.0"
	AppDescription = "A desktop application for managing and downloading lyrics"
	AppAuthor      = "LRCGET Team"
	AppURL         = "https://github.com/tranxuanthang/lrcget"
)

// Database constants
const (
	DatabaseVersion  = 7
	DatabaseFileName = "db.sqlite3"
	DefaultDataDir   = "~/.lrcget"
	MaxDatabaseSize  = 100 * 1024 * 1024 // 100MB
	DatabaseTimeout  = 30 * time.Second
	MaxRetries       = 3
	RetryDelay       = 1 * time.Second
)

// HTTP client constants
const (
	DefaultTimeout      = 30 * time.Second
	DefaultUserAgent    = "LRCGET v1.0.0 (https://github.com/tranxuanthang/lrcget)"
	MaxRequestSize      = 10 * 1024 * 1024 // 10MB
	MaxResponseSize     = 50 * 1024 * 1024 // 50MB
	MaxRedirects        = 5
	KeepAliveTimeout    = 30 * time.Second
	IdleConnTimeout     = 90 * time.Second
	MaxIdleConns        = 10
	MaxIdleConnsPerHost = 2
)

// Worker pool constants
const (
	DefaultMaxWorkers = 10
	MinWorkers        = 1
	MaxWorkers        = 100
	WorkerTimeout     = 5 * time.Minute
	JobQueueSize      = 100
	ResultQueueSize   = 100
)

// File system constants
const (
	MaxFileSizeForMemoryProcessing = 10 * 1024 * 1024  // 10MB
	MaxFileSizeForStreaming        = 100 * 1024 * 1024 // 100MB
	DefaultBufferSize              = 4096
	MaxBufferSize                  = 64 * 1024 // 64KB
	MinBufferSize                  = 1024      // 1KB
)

// Audio file constants
const (
	MaxAudioFileSize = 500 * 1024 * 1024 // 500MB
	MinAudioFileSize = 1024              // 1KB
	MaxDuration      = 3600              // 1 hour in seconds
	MinDuration      = 1                 // 1 second
)

// Lyrics constants
const (
	MaxLyricsLength      = 100000 // 100KB
	MinLyricsLength      = 10
	MaxSearchResults     = 100
	DefaultSearchLimit   = 20
	MaxSearchQueryLength = 1000
	MinSearchQueryLength = 1
)

// Cache constants
const (
	DefaultCacheSize     = 100
	MaxCacheSize         = 10000
	CacheExpiration      = 1 * time.Hour
	CacheCleanupInterval = 5 * time.Minute
	MaxCacheItemSize     = 1024 * 1024 // 1MB
)

// Logging constants
const (
	DefaultLogLevel     = "info"
	DefaultLogFormat    = "json"
	DefaultLogOutput    = "stdout"
	MaxLogFileSize      = 100 * 1024 * 1024 // 100MB
	MaxLogFiles         = 3
	LogRotationInterval = 24 * time.Hour
)

// Security constants
const (
	MinPasswordLength    = 8
	MaxPasswordLength    = 128
	MinUsernameLength    = 3
	MaxUsernameLength    = 50
	MaxLoginAttempts     = 5
	LoginLockoutDuration = 15 * time.Minute
	SessionTimeout       = 24 * time.Hour
	TokenExpiration      = 1 * time.Hour
)

// API constants
const (
	DefaultLRCLibURL      = "https://lrclib.net"
	APIVersion            = "v1"
	MaxConcurrentRequests = 10
	RequestTimeout        = 30 * time.Second
	RetryAttempts         = 3
)

// UI constants
const (
	DefaultWindowWidth  = 1200
	DefaultWindowHeight = 800
	MinWindowWidth      = 800
	MinWindowHeight     = 600
	DefaultTheme        = "light"
	DefaultFontSize     = 14
	MinFontSize         = 10
	MaxFontSize         = 24
)

// Player state constants
const (
	StatusStopped = iota
	StatusPlaying
	StatusPaused
	StatusBuffering
	StatusErrorState
)

// File type constants
const (
	AudioFileType  = "audio"
	LyricsFileType = "lyrics"
	ConfigFileType = "config"
	LogFileType    = "log"
	CacheFileType  = "cache"
)

// Audio format constants
const (
	MP3Format  = "mp3"
	M4AFormat  = "m4a"
	FLACFormat = "flac"
	OGGFormat  = "ogg"
	OPUSFormat = "opus"
	WAVFormat  = "wav"
	AACFormat  = "aac"
	WMAFormat  = "wma"
)

// Lyrics format constants
const (
	LRCFormat = "lrc"
	TXTFormat = "txt"
	SRTFormat = "srt"
	VTTFormat = "vtt"
)

// Error codes
const (
	ErrCodeInvalidInput        = "INVALID_INPUT"
	ErrCodeFileNotFound        = "FILE_NOT_FOUND"
	ErrCodePermissionDenied    = "PERMISSION_DENIED"
	ErrCodeNetworkError        = "NETWORK_ERROR"
	ErrCodeDatabaseError       = "DATABASE_ERROR"
	ErrCodeValidationError     = "VALIDATION_ERROR"
	ErrCodeAuthenticationError = "AUTHENTICATION_ERROR"
	ErrCodeAuthorizationError  = "AUTHORIZATION_ERROR"
	ErrCodeRateLimitError      = "RATE_LIMIT_ERROR"
	ErrCodeServerError         = "SERVER_ERROR"
	ErrCodeUnknownError        = "UNKNOWN_ERROR"
)

// Status codes
const (
	StatusSuccess = "success"
	StatusError   = "error"
	StatusWarning = "warning"
	StatusInfo    = "info"
)

// Environment variables
const (
	EnvDatabasePath  = "LRCGET_DB_PATH"
	EnvMaxWorkers    = "LRCGET_MAX_WORKERS"
	EnvTimeout       = "LRCGET_TIMEOUT"
	EnvLRCLibURL     = "LRCGET_LRCLIB_URL"
	EnvLogLevel      = "LRCGET_LOG_LEVEL"
	EnvDataDir       = "LRCGET_DATA_DIR"
	EnvCacheSize     = "LRCGET_CACHE_SIZE"
	EnvEnableMetrics = "LRCGET_ENABLE_METRICS"
	EnvDebugMode     = "LRCGET_DEBUG"
	EnvConfigFile    = "LRCGET_CONFIG_FILE"
)

// Default values
const (
	DefaultConfigFile = "config.json"
	DefaultLogFile    = "lrcget.log"
	DefaultCacheDir   = "cache"
	DefaultTempDir    = "temp"
	DefaultBackupDir  = "backup"
)

// Validation constants
const (
	MaxPathLength     = 4096
	MaxFilenameLength = 255
	MaxDirectoryDepth = 100
	MaxFileCount      = 100000
	MaxDirectoryCount = 10000
)

// Performance constants
const (
	DefaultScanInterval        = 5 * time.Minute
	DefaultCleanupInterval     = 1 * time.Hour
	DefaultMetricsInterval     = 1 * time.Minute
	DefaultHealthCheckInterval = 30 * time.Second
)

// Quality constants
const (
	QualityLow    = 1
	QualityMedium = 2
	QualityHigh   = 3
	QualityBest   = 4
)

// Priority constants
const (
	PriorityLow      = 1
	PriorityMedium   = 2
	PriorityHigh     = 3
	PriorityCritical = 4
)
