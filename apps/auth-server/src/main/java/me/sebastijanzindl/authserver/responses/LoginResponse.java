package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import me.sebastijanzindl.authserver.dto.UserDTO;

import java.util.Date;

public record LoginResponse(
        @JsonProperty("user") UserDTO user,
        @JsonProperty("access_token") String accessToken,
        @JsonProperty("refresh_token") String refreshToken,
        @JsonProperty("access_token_expires_in") Date accessTokenExpiration,
        @JsonProperty("refresh_token_expires_in") Date refreshTokenExpiration)
{
}
