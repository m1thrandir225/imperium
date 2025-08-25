package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
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
    @NotNull(message = "Port cannot be null")
    @Min(1)
    @Max(65535)
    private Integer port;
}
