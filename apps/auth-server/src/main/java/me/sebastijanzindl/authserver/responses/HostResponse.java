package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.Getter;
import me.sebastijanzindl.authserver.model.Host;

import java.util.UUID;

public record HostResponse (
    @JsonProperty("id")
     UUID id,
    @JsonProperty("ip_address")
     String ipAddress,
    @JsonProperty("port")
     Integer port,
    @JsonProperty("name")
     String name,
    @JsonProperty("status")
     String status
){}
