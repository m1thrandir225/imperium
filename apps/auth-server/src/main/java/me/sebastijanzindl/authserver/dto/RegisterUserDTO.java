package me.sebastijanzindl.authserver.dto;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class RegisterUserDTO {
    private String email;
    private String password;
    private String firstName;
    private String lastName;
}
