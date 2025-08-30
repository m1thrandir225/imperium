import configService from "@/services/config.service";
import type {ConfigStore} from "@/types/stores/config";
import {create} from "zustand";
import {persist} from "zustand/middleware";

const useConfigStore = create<ConfigStore>()(
  persist(
    (set) => ({
      isConfigured: false,
      isLoading: true,
      error: null,

      checkConfiguration: async () => {
        try {
          set({isLoading: true, error: null});
          const status = await configService.getConfigStatus();
          set({isConfigured: status.configured, isLoading: false});
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : "Unknown error",
            isLoading: false,
          });
        }
      },

      setupConfiguration: (configured: boolean) => {
        set({isConfigured: configured});
      },
    }),
    {
      name: "config",
      partialize: (state) => ({
        isConfigured: state.isConfigured,
      }),
    }
  )
);

export default useConfigStore;
