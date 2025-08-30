import type {User} from "../models/user";

/*
  This is the request to login.
*/
export type LoginRequest = {
  email: string;
  password: string;
};

/*
  This is the response to login.
*/
export type LoginResponse = {
  user: User;
  access_token: string;
  refresh_token: string;
  access_token_expires_in: string;
  refresh_token_expires_in: string;
};

/*
  This is the request to register.
*/
export type RegisterRequest = {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
};

/*
  This is the response to register.
*/
export type RegisterResponse = {
  user: User;
};

/*
  This is the request to refresh the access token.
*/
export type RefreshTokenRequest = {
  token: string;
};

/*
  This is the response to refresh the access token.
*/
export type RefreshTokenResponse = {
  access_token: string;
  expires_at: string;
};
