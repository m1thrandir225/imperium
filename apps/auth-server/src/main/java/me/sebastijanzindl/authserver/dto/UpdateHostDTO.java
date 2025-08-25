package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.Getter;
import me.sebastijanzindl.authserver.model.enums.HOST_STATUS;

@Data
@AllArgsConstructor
public class UpdateHostDTO {
    @JsonProperty("ip_address")
    @NotBlank(message = "IP address cannot be blank")
    @NotNull(message = "IP address cannot be null")
    private String ipAddress;

    @JsonProperty("port")
    @NotNull(message = "Port cannot be null")
    @NotBlank(message = "Port cannot be blank")
    private Integer port;

    @JsonProperty("status")
    @NotNull(message = "Status cannot be null")
    @NotBlank(message = "Status cannot be blank")
    private HOST_STATUS status;
}
