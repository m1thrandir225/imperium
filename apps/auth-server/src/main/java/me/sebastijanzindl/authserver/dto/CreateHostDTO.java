package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.Getter;

@Getter
public class CreateHostDTO {
    @JsonProperty("name")
    @NotBlank(message = "Name cannot be blank")
    @NotNull(message = "Name cannot be null")
    private String name;


    @JsonProperty("ip_address")
    @NotBlank(message = "IP Address cannot be blank")
    @NotNull(message = "IP Address cannot be null")
    private String ipAddress;

    @JsonProperty("port")
    @NotBlank(message = "Port cannot be blank")
    @NotNull(message = "Port cannot be null")
    private Integer port;
}
