type State = {
  isConfigured: boolean;
  isLoading: boolean;
  error: string | null;
};

type Actions = {
  checkConfiguration: () => Promise<void>;
  setupConfiguration: (configured: boolean) => void;
};

export type ConfigStore = State & Actions;
