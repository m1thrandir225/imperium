import type {Host} from "../models/host";
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
  host: Host | null;
};

type Actions = {
  setHost: (host: Host) => void;
  setCurrentSession: (session: Session) => void;
  createSession: (request: CreateSessionRequest) => Promise<Session>;
  startSession: (
    request: StartSessionRequest,
    sessionId: string
  ) => Promise<Session>;

  clearError: () => void;
};
