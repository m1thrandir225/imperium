package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotNull;
import lombok.Data;

import java.util.UUID;

@Data
public class CreateSessionDTO {
    @JsonProperty("host_id")
    @NotNull(message = "Host ID is required")
    private UUID hostId;

    @JsonProperty("client_id")
    @NotNull(message = "Client ID is required")
    private UUID clientId;

    @JsonProperty("program_id")
    private String programId;
}
