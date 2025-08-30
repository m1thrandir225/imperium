package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.Getter;
import me.sebastijanzindl.authserver.model.Host;

import java.util.List;
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
){
    public static HostResponse from(Host host) {
        return new HostResponse(
                host.getId(),
                host.getIpAddress(),
                host.getPort(),
                host.getName(),
                host.getStatus().name()
        );
    }
    public static List<HostResponse> fromList(List<Host> hosts) {
        return hosts.stream()
                .map(HostResponse::from)
                .toList();
    }
}
