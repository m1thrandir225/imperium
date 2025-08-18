package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class LoginUserDTO {
    @JsonProperty("email")
    private String email;

    @JsonProperty("password")
    private String password;
}
