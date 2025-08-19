package me.sebastijanzindl.authserver.dto;

import lombok.Getter;

@Getter
public class CreateHostDTO {
    private String name;
    private String ipAddress;
    private Integer port;
}
