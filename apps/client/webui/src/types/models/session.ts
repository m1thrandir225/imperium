export type SessionStatus =
  | "PENDING"
  | "ACTIVE"
  | "ENDED"
  | "CANCELLED"
  | "FAILED";

export type Session = {
  id: string;
  host_id: string;
  host_name: string;
  client_id: string;
  client_name: string;
  status: SessionStatus;
  session_token: string;
  webrtc_offer?: string;
  webrtc_answer?: string;
  expires_at: string;
  created_at: string;
  started_at?: string;
  ended_at?: string;
  end_reason?: string;
};

export type CreateSessionRequest = {
  host_id: string;
  client_id: string;
  program_id: string;
};

export type StartSessionRequest = {
  webrtc_offer: string;
};

export type EndSessionRequest = {
  reason?: string;
  webrtc_answer?: string;
};

export type CancelSessionRequest = {
  reason: string;
};
