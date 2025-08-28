import type {AuthStore} from "@/types/stores/auth";
import {create} from "zustand";
import {persist, createJSONStorage} from "zustand/middleware";

/*
  Persisted auth store, used for token and user data.
*/
const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      accessToken: null,
      accessTokenExpiresAt: null,
      refreshToken: null,
      refreshTokenExpiresAt: null,
      user: null,
      isHydrated: false,
      setUser: (newUser) => {
        set({user: newUser});
      },
      login: (data) => {
        set({
          user: data.user,
          accessToken: data.access_token,
          accessTokenExpiresAt: new Date(data.access_token_expires_in),
          refreshToken: data.refresh_token,
          refreshTokenExpiresAt: new Date(data.refresh_token_expires_in),
        });
      },
      refreshAccessToken: (data) => {
        set({
          accessToken: data.access_token,
          accessTokenExpiresAt: new Date(data.expires_at),
        });
      },
      checkAuth: (retryNumber) => {
        if (retryNumber === undefined) {
          retryNumber = 0;
        }

        if (retryNumber && retryNumber > 5) {
          get().logout();
          return false;
        }

        const isHydrated = get().isHydrated;
        if (!isHydrated) {
          setTimeout(() => {
            get().checkAuth(retryNumber + 1);
          }, 1000);
          return false;
        }
        if (!get().canRefresh()) {
          get().logout();
          return false;
        }
        return true;
      },
      canRefresh: () => {
        const now = new Date();
        const refreshTokenExpiresAt = get().refreshTokenExpiresAt;
        const refreshToken = get().refreshToken;
        if (
          !refreshToken ||
          !refreshTokenExpiresAt ||
          now > refreshTokenExpiresAt
        ) {
          return false;
        }
        return true;
      },
      logout: () => {
        set({
          user: null,
          accessToken: null,
          accessTokenExpiresAt: null,
          refreshToken: null,
          refreshTokenExpiresAt: null,
        });
      },
      setHasHydrated: (isHydrated) => {
        set({isHydrated});
      },
    }),
    {
      name: "auth",
      storage: createJSONStorage(() => localStorage),
      onRehydrateStorage: (state) => {
        return () => state.setHasHydrated(true);
      },
    }
  )
);

export default useAuthStore;
