package me.sebastijanzindl.authserver.dto;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class UpdateUserPasswordDTO {
    private String password;
    private String newPassword;

    public UpdateUserPasswordDTO(String password, String newPassword) {
        this.newPassword = newPassword;
        this.password = password;
    }
}
