import config from "@/lib/config";
import {
  type Session,
  type CreateSessionRequest,
  type EndSessionRequest,
  type StartSessionRequest,
} from "@/types/models/session";
import {apiRequest} from "./api.service";

const sessionServiceAPIUrl = `${config.apiUrl}/sessions`;

const sessionService = {
  create: (input: CreateSessionRequest) =>
    apiRequest<Session>({
      url: `${sessionServiceAPIUrl}/create`,
      method: "POST",
      data: input,
      protected: true,
      headers: undefined,
      params: undefined,
    }),
  start: (input: StartSessionRequest, sessionId: string) =>
    apiRequest<Session>({
      url: `${sessionServiceAPIUrl}/${sessionId}/start`,
      method: "POST",
      data: input,
      protected: true,
      headers: undefined,
      params: undefined,
    }),
  end: (input: EndSessionRequest, sessionId: string) =>
    apiRequest<Session>({
      url: `${sessionServiceAPIUrl}/${sessionId}/end`,
      method: "POST",
      data: input,
      protected: true,
      headers: undefined,
      params: undefined,
    }),
  get: (sessionId: string) =>
    apiRequest<Session>({
      url: `${sessionServiceAPIUrl}/${sessionId}`,
      method: "GET",
      protected: true,
      headers: undefined,
      params: undefined,
    }),
  pollUntilActive: async (
    sessionId: string,
    onUpdate?: (s: Session) => void,
    maxAttempts = 30,
    interval = 2000
  ) => {
    for (let i = 0; i < maxAttempts; i++) {
      const session = await sessionService.get(sessionId);
      onUpdate?.(session);

      if (session.status === "ACTIVE") return session;
      if (
        session.status === "FAILED" ||
        session.status === "CANCELLED" ||
        session.status === "ENDED"
      ) {
        throw new Error(
          `Session ${session.status.toLowerCase()}${
            session.end_reason ? `: ${session.end_reason}` : ""
          }`
        );
      }

      await new Promise((r) => setTimeout(r, interval));
    }
    throw new Error("Session did not become active in time");
  },
};

export default sessionService;
