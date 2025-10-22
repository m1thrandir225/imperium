package programs

import "errors"

var (
	ErrDatabaseDirectoryCreationFailed = errors.New("failed to create directory")
	ErrDatabaseOpenFailed              = errors.New("failed to open database")
	ErrDatabasePingFailed              = errors.New("failed to ping database")
	ErrDatabaseInitializeTablesFailed  = errors.New("failed to initialize tables")
	ErrOSNotSupported                  = errors.New("this OS is currently not supported")
)
