package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import me.sebastijanzindl.authserver.model.Session;
import me.sebastijanzindl.authserver.model.enums.SESSION_STATUS;

import java.time.LocalDateTime;
import java.util.UUID;

public record SimpleSessionResponse(
        @JsonProperty("id")
        UUID id,

        @JsonProperty("host_name")
        String hostName,

        @JsonProperty("client_name")
        String clientName,

        @JsonProperty("status")
        SESSION_STATUS status,

        @JsonProperty("created_at")
        LocalDateTime createdAt,

        @JsonProperty("started_at")
        LocalDateTime startedAt,

        @JsonProperty("ended_at")
        LocalDateTime endedAt
) {
    public SimpleSessionResponse(Session session) {
        this(
                session.getId(),
                session.getHost().getName(),
                session.getClient().getName(),
                session.getStatus(),
                session.getCreatedAt(),
                session.getStartedAt(),
                session.getEndedAt()
        );
    }
}
