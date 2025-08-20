package host

type Config struct {
	HostName  string `json:"host_name"`
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	UniqueID  string `json:"unique_id"`
	Status    string `json:"status"`
}

func NewConfig(hostname, ipAddress string, port int, uniqueID string) *Config {
	return &Config{
		HostName:  hostname,
		IPAddress: ipAddress,
		Port:      port,
		UniqueID:  uniqueID,
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
