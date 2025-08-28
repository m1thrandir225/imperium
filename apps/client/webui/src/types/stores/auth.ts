import type {User} from "../models/user";
import type {LoginResponse, RefreshTokenResponse} from "../responses/auth";

export type AuthStore = State & Actions;

type State = {
  accessToken: string | null;
  accessTokenExpiresAt: Date | null;
  refreshToken: string | null;
  refreshTokenExpiresAt: Date | null;
  user: User | null;

  isHydrated: boolean;
};

type Actions = {
  setUser: (user: User) => void;
  login: (data: LoginResponse) => void;
  refreshAccessToken: (data: RefreshTokenResponse) => void;
  canRefresh: () => boolean;
  checkAuth: (retryNumber?: number) => boolean;
  logout: () => void;
  setHasHydrated: (isHydrated: boolean) => void;
};
