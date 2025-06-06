package me.sebastijanzindl.authserver.dto;

import lombok.Getter;

@Getter
public class CreateClientDTO {
    private String ipAddress;
    private String clientName;
}
