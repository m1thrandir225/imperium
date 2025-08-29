package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import me.sebastijanzindl.authserver.model.enums.SESSION_STATUS;

import java.time.LocalDateTime;
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
        String webrtcOffer,

        @JsonProperty("answer")
        String webrtcAnswer,

        @JsonProperty("expires_at")
        LocalDateTime expiresAt,

        @JsonProperty("created_at")
        LocalDateTime createdAt,

        @JsonProperty("updated_at")
        LocalDateTime startedAt,

        @JsonProperty("ended_at")
        LocalDateTime endedAt,

        @JsonProperty("end_reason")
        String endReason
) {
}
