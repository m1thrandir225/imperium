import type {
  CreateSessionRequest,
  Session,
  StartSessionRequest,
} from "../models/session";

export type SessionStore = State & Actions;

type State = {
  currentSession: Session | null;
  isConnecting: boolean;
  connectionError: string | null;
};

type Actions = {
  setCurrentSession: (session: Session) => void;
  createSession: (request: CreateSessionRequest) => Promise<Session>;
  startSession: (
    request: StartSessionRequest,
    sessionId: string
  ) => Promise<Session>;

  clearError: () => void;
};
