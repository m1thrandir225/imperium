package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;

@Getter
public class CreateHostDTO {
    @JsonProperty("name")
    private String name;

    @JsonProperty("ip_address")
    private String ipAddress;

    @JsonProperty("port")
    private Integer port;
}
