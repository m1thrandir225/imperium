// Package host provides the functionality to manage the host machine.
// This includes the functionality to get the host machine's unique ID,
// Get the host machine's IP address,
// Get the host running port,
package host

type Host struct {
	UniqueID  string
	IPAddress string
	Port      string
	HostName  string
	OS        string
	Status    string
}
