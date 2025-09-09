package host

type Status string

const (
	StatusAvailable Status = "available"
	StatusOffline   Status = "offline"
	StatusInuse     Status = "inuse"
	StatusDisabled  Status = "disabled"
	StatusUnknown   Status = "unknown"
)

func (s Status) toAPIEnum() string {
	switch s {
	case StatusAvailable:
		return "AVAILABLE"
	case StatusOffline:
		return "OFFLINE"
	case StatusInuse:
		return "INUSE"
	case StatusDisabled:
		return "DISABLED"
	case StatusUnknown:
		return "UNKNOWN"
	default:
		return "UNKNOWN"

	}
}
