package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import me.sebastijanzindl.authserver.model.Host;

import java.util.UUID;

@Getter
public class HostResponse {
    @JsonProperty("id")
    private UUID id;

    @JsonProperty("ip_address")
    private String ipAddress;

    @JsonProperty("port")
    private Integer port;

    @JsonProperty("name")
    private String name;

    @JsonProperty("status")
    private String status;

    public HostResponse(Host host) {
        this.id = host.getId();
        this.ipAddress = host.getIpAddress();
        this.port = host.getPort();
        this.name = host.getName();
        this.status = host.getStatus().name();
    }
}
