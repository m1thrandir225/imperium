package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import me.sebastijanzindl.authserver.model.Session;
import me.sebastijanzindl.authserver.model.enums.SESSION_STATUS;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

public record SessionResponse(
        @JsonProperty("id")
        UUID id,

        @JsonProperty("host_id")
        UUID hostId,

        @JsonProperty("host_name")
        String hostName,

        @JsonProperty("client_id")
        UUID clientId,

        @JsonProperty("client_name")
        String clientName,

        @JsonProperty("status")
        SESSION_STATUS status,

        @JsonProperty("session_token")
        String sessionToken,

        @JsonProperty("program_id")
        String programId,

        @JsonProperty("webrtc_offer")
        String webrtcOffer,

        @JsonProperty("webrtc_answer")
        String webrtcAnswer,

        @JsonProperty("expires_at")
        LocalDateTime expiresAt,

        @JsonProperty("created_at")
        LocalDateTime createdAt,

        @JsonProperty("started_at")
        LocalDateTime startedAt,

        @JsonProperty("ended_at")
        LocalDateTime endedAt,

        @JsonProperty("end_reason")
        String endReason
) {
    public static SessionResponse from(Session session) {
        return new SessionResponse(
                session.getId(),
                session.getHost().getId(),
                session.getHost().getName(),
                session.getClient().getId(),
                session.getClient().getName(),
                session.getStatus(),
                session.getSessionToken(),
                session.getProgramId(),
                session.getWebrtcOffer(),
                session.getWebrtcAnswer(),
                session.getExpiresAt(),
                session.getCreatedAt(),
                session.getStartedAt(),
                session.getEndedAt(),
                session.getEndReason()

        );
    }

    public static List<SessionResponse> fromList(List<Session> sessions) {
        return sessions.stream()
                .map(SessionResponse::from)
                .toList();
    }
}
