package me.sebastijanzindl.authserver.dto;

import lombok.Getter;
import me.sebastijanzindl.authserver.model.enums.HOST_STATUS;

@Getter
public class UpdateHostDTO {
    private String ipAddress;
    private Integer port;
    private HOST_STATUS status;
}
