import sessionService from "@/services/session.service";
import type {SessionStore} from "@/types/stores/session";
import {create} from "zustand";

export const useSessionStore = create<SessionStore>()((set) => ({
  currentSession: null,
  isConnecting: false,
  connectionError: null,
  host: null,

  setCurrentSession: (session) => {
    set({currentSession: session});
  },
  setHost: (host) => {
    set({host: host});
  },
  createSession: async (request) => {
    set({isConnecting: true, connectionError: null});
    try {
      const s = await sessionService.create(request);
      set({currentSession: s});
      return s;
    } catch (e: unknown) {
      set({connectionError: e instanceof Error ? e.message : "Unknown error"});
      throw e;
    } finally {
      set({isConnecting: false});
    }
  },

  startSession: async (request, sessionId) => {
    set({isConnecting: true, connectionError: null});
    try {
      const s = await sessionService.start(request, sessionId);
      set({currentSession: s});
      return s;
    } catch (e: unknown) {
      set({connectionError: e instanceof Error ? e.message : "Unknown error"});
      throw e;
    } finally {
      set({isConnecting: false});
    }
  },

  clearError: () => {
    set({connectionError: null});
  },
}));
