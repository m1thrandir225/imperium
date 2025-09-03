package host

import "github.com/m1thrandir225/imperium/apps/host/internal/util"

type Config struct {
	HostName  string `json:"host_name"`
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	UniqueID  string `json:"unique_id"`
	Status    string `json:"status"`
}

func NewConfig(uniqueID string) *Config {
	hostname, err := util.GetHostname()
	if err != nil {
		hostname = ""
	}

	ipAddress, err := util.GetIPAddress()
	if err != nil {
		ipAddress = ""
	}

	port := 8080

	return &Config{
		HostName:  hostname,
		IPAddress: ipAddress,
		Port:      port,
		UniqueID:  uniqueID,
		Status:    string(StatusAvailable),
	}
}

func (c *Config) GetHostName() string {
	return c.HostName
}

func (c *Config) GetIPAddress() string {
	return c.IPAddress
}

func (c *Config) GetPort() int {
	return c.Port
}

func (c *Config) GetUniqueID() string {
	return c.UniqueID
}

func (c *Config) GetStatus() string {
	return c.Status
}
