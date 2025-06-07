package me.sebastijanzindl.authserver.dto;

import lombok.Getter;

@Getter
public class UpdateUserDTO {
    private String email;
    private String name;
}
