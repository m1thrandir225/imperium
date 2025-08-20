package host

type Status string

const (
	StatusAvailable Status = "available"
	StatusOffline   Status = "offline"
	StatusInuse     Status = "inuse"
	StatusDisabled  Status = "disabled"
	StatusUnknown   Status = "unknown"
)
