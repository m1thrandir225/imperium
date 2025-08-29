package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import me.sebastijanzindl.authserver.model.enums.SESSION_STATUS;

import java.time.LocalDateTime;
import java.util.UUID;

public record SimpleSessionResponse(
        @JsonProperty("id")
        UUID id,

        @JsonProperty("host_id")
        String hostName,

        @JsonProperty("client_id")
        String clientName,

        @JsonProperty("status")
        SESSION_STATUS status,

        @JsonProperty("created_at")
        LocalDateTime createdAt,

        @JsonProperty("updated_at")
        LocalDateTime startedAt,

        @JsonProperty("ended_at")
        LocalDateTime endedAt
) {
}
