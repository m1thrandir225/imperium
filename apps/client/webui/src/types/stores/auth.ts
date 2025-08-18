export type AuthStore = State & Actions;

type State = {
  accessToken: string | null;
  refreshToken: string | null;
};

type Actions = {
  setTokens: (accessToken: string | null, refreshToken: string | null) => void;
  checkAuth: () => boolean;
  logout: () => void;
};
