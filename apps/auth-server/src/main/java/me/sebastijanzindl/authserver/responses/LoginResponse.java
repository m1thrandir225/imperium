package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import lombok.Setter;
import me.sebastijanzindl.authserver.dto.LoginUserDTO;
import me.sebastijanzindl.authserver.dto.UserDTO;

import java.util.Date;

@Getter
@Setter
public class LoginResponse {
    @JsonProperty("user")
    private UserDTO user;

    @JsonProperty("access_token")
    private String accessToken;

    @JsonProperty("refresh_token")
    private String refreshToken;

    @JsonProperty("access_token_expires_in")
    private Date accessTokenExpiration;

    @JsonProperty("refresh_token_expires_in")
    private Date refreshTokenExpiration;
}
