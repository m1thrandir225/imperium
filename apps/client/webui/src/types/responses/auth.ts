export type LoginRequest = {
  email: string;
  password: string;
};

export type LoginResponse = {
  token: string;
  refreshToken: string;
  expiresIn: number;
};

export type RegisterRequest = {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
};

export type RegisterResponse = {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  createdAt: string;
  updatedAt: string;
  clients: unknown[];
  hosts: unknown[];
};

export type RefreshTokenRequest = {
  refresh_token: string;
};

export type RefreshTokenResponse = {
  newToken: string;
};
