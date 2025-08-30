export type ConfigStatus = {
  configured: boolean;
};

export type SetupConfigRequest = {
  auth_server_base_url: string;
};

export type SetupConfigResponse = {
  message: string;
};
