package me.sebastijanzindl.authserver.responses;

import lombok.Getter;
import lombok.Setter;
import me.sebastijanzindl.authserver.dto.LoginUserDTO;

@Getter
@Setter
public class LoginResponse {
    private String token;
    private String refreshToken;
    private Long expiresIn;
}
