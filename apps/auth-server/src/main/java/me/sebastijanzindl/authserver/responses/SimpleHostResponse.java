package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import me.sebastijanzindl.authserver.model.Host;

import java.util.List;
import java.util.UUID;


public record SimpleHostResponse(
        @JsonProperty("id")
        UUID id,

        @JsonProperty("name")
        String name,

        @JsonProperty("status")
        String status
) {
    public static SimpleHostResponse from(Host host) {
        return new SimpleHostResponse(
                host.getId(),
                host.getName(),
                host.getStatus().name()
        );
    }

    public static List<SimpleHostResponse> fromList(List<Host> hosts) {
        return hosts.stream()
                .map(SimpleHostResponse::from)
                .toList();
    }
}
