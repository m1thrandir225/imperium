package me.sebastijanzindl.authserver.dto;

import lombok.Getter;

@Getter
public class CreateHostDTO {
    private String ipAddress;
    private Integer port;
}
