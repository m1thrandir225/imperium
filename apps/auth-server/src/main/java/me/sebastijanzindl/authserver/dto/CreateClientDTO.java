package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.Getter;

@Getter
public class CreateClientDTO {
    @JsonProperty("ip_address")
    @NotNull(message = "IP address cannot be null")
    @NotBlank(message = "IP address cannot be blank")
    private String ipAddress;

    @JsonProperty("client_name")
    @NotNull(message = "Client name cannot be null")
    @NotBlank(message = "Client name cannot be blank")
    private String clientName;
}