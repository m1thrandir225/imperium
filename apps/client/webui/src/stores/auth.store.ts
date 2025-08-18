import type {AuthStore} from "@/types/stores/auth";
import {create} from "zustand";

const useAuthStore = create<AuthStore>((set, get) => ({
  accessToken: null,
  refreshToken: null,
  setTokens: (accessToken, refreshToken) => set({accessToken, refreshToken}),
  checkAuth: () => {
    const {accessToken, refreshToken} = get();
    return !!(accessToken && refreshToken);
  },
  logout: () => {
    set({accessToken: null, refreshToken: null});
  },
}));

export default useAuthStore;
