package system

import (
	"errors"
	"fmt"
)

var (
	ErrExec = func(e error) error {
		return fmt.Errorf("system service execution error: %s", e.Error())
	}
)

var (
	ExitCodesMap = map[int]string{
		1:   "Command execution failed",
		2:   "Invalid arguments provided",
		3:   "Permission denied",
		4:   "Memory overflow error",
		5:   "Input/Output error",
		6:   "File not found",
		7:   "Invalid file descriptor",
		8:   "Memory allocation error",
		9:   "Binary execution error",
		10:  "Environment variable error",
		11:  "Interrupted by signal (e.g., SIGINT)",
		12:  "Resource limit exceeded",
		13:  "Network error",
		14:  "Operation timed out",
		15:  "Authentication error",
		16:  "Configuration error",
		17:  "Dependency error",
		18:  "Unknown error occurred",
		19:  "Lock acquisition error",
		20:  "Database error",
		21:  "Cache error",
		22:  "Syntax error",
		23:  "Invalid state",
		24:  "Duplicate entry",
		25:  "Conflict detected",
		26:  "Value out of range",
		27:  "Invalid format",
		28:  "Invalid operation",
		29:  "Operation not supported",
		30:  "Internal error",
		31:  "Deployment error",
		32:  "Parsing error",
		33:  "File system error",
		34:  "Stream error",
		35:  "Thread-related error",
		36:  "Module error",
		37:  "Process interrupted",
		38:  "Operation failed",
		39:  "Retry limit exceeded",
		40:  "Compatibility issue",
		41:  "Data corruption detected",
		42:  "Resource not found",
		43:  "Invalid input provided",
		44:  "Illegal state detected",
		45:  "Service unavailable",
		46:  "Interrupt error",
		47:  "Unsupported file type",
		48:  "Illegal operation",
		49:  "Deadlock detected",
		50:  "Process limit exceeded",
		51:  "Terminated by signal",
		52:  "Stack overflow error",
		53:  "Initialization error",
		54:  "Invalid response received",
		55:  "API error",
		56:  "Deprecated feature used",
		57:  "Version mismatch",
		58:  "Invalid permissions",
		59:  "Hardware-related error",
		60:  "Unhandled exception",
		61:  "Logic error in code",
		62:  "Concurrency error",
		63:  "Security violation detected",
		64:  "Invalid token provided",
		65:  "Numeric overflow error",
		66:  "Numeric underflow error",
		67:  "Protocol error",
		68:  "Illegal character found",
		69:  "Session error",
		70:  "Transaction error",
		71:  "Validation failed",
		72:  "File lock error",
		73:  "Connection reset by peer",
		74:  "Buffer overflow error",
		75:  "Buffer underflow error",
		76:  "Encoding error",
		77:  "Decoding error",
		78:  "Quota exceeded",
		79:  "Key not found or invalid",
		80:  "Invalid value provided",
		81:  "Redirection error",
		82:  "Shutdown error",
		83:  "Invalid path",
		84:  "Feature not implemented",
		85:  "General service error",
		86:  "Process error",
		87:  "Blocked by policy",
		88:  "License error",
		89:  "Connection error",
		90:  "Disk error",
		91:  "Unsupported media type",
		92:  "Policy violation",
		93:  "Restriction error",
		94:  "I/O operation timed out",
		95:  "Resource is locked",
		96:  "Upgrade error",
		97:  "Migration error",
		98:  "DNS resolution error",
		99:  "Token expired",
		100: "Certificate error",
		101: "Redis-related error",
		102: "Runtime error",
		103: "Command-line interface error",
		104: "Unknown signal received",
		105: "Monitoring failure",
		106: "Heartbeat missed",
		107: "Module not found",
		108: "Operation not allowed",
		109: "Job execution error",
		110: "Scheduling error",
		111: "External service error",
		112: "Function execution error",
		113: "Test failure",
		114: "Parallelism error",
		115: "Analytics error",
		116: "Script execution error",
		117: "Session expired",
		118: "Workflow error",
		119: "Dependency timed out",
		120: "Cluster operation error",
	}
)

type ExitErrorCode struct {
	Code        int
	Description string
}

func NewExitErrorCode(code int) *ExitErrorCode {

	res, ok := ExitCodesMap[code]

	if !ok {
		return nil
	}

	return &ExitErrorCode{
		Code:        code,
		Description: res,
	}
}

func (e *ExitErrorCode) Error() string {
	return fmt.Sprintf("system service execution error: %s", e.Description)
}

var (
	ErrUnitNotFound = errors.New("systemd unit not found")
)
