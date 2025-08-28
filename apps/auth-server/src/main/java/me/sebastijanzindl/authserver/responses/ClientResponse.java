package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import me.sebastijanzindl.authserver.model.Client;

import java.util.UUID;

public record ClientResponse (
    @JsonProperty("id")
     UUID id,
    @JsonProperty("client_name")
     String clientName,
    @JsonProperty("ip_address")
     String ipAddress
)
{
}
