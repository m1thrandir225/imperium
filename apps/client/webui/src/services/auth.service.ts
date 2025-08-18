import type {
  LoginRequest,
  LoginResponse,
  RefreshTokenRequest,
  RefreshTokenResponse,
  RegisterRequest,
  RegisterResponse,
} from "@/types/responses/auth";
import {apiRequest} from "./api.service";
import config from "@/lib/config";
const authURL = config.apiUrl;

const authService = {
  login: (input: LoginRequest) =>
    apiRequest<LoginResponse>({
      method: "POST",
      url: `${authURL}/login`,
      data: input,
      protected: false,
      headers: undefined,
      params: undefined,
    }),
  register: (input: RegisterRequest) =>
    apiRequest<RegisterResponse>({
      method: "POST",
      url: `${authURL}/register`,
      data: input,
      protected: false,
      headers: undefined,
      params: undefined,
    }),
  refreshToken: (input: RefreshTokenRequest) =>
    apiRequest<RefreshTokenResponse>({
      method: "POST",
      data: input,
      protected: false,
      headers: undefined,
      params: undefined,
      url: `${authURL}/refresh-token`,
    }),
};

export default authService;
