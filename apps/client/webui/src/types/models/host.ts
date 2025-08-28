export type Host = {
  id: string;
  ip_address: string;
  port: number;
  name: string;
  status: "OFFLINE" | "AVAILABLE" | "INUSE" | "DISABLED";
};

export type SimpleHost = Omit<Host, "ip_address" | "port">;
