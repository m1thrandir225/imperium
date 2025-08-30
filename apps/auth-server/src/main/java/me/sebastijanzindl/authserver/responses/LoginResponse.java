package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import me.sebastijanzindl.authserver.dto.UserDTO;
import me.sebastijanzindl.authserver.model.User;

import java.util.Date;

public record LoginResponse(
        @JsonProperty("user") UserDTO user,
        @JsonProperty("access_token") String accessToken,
        @JsonProperty("refresh_token") String refreshToken,
        @JsonProperty("access_token_expires_in") Date accessTokenExpiration,
        @JsonProperty("refresh_token_expires_in") Date refreshTokenExpiration)
{
    public static LoginResponse from(
            User user,
            String accessToken,
            String refreshToken,
            Date accessTokenExpiration,
            Date refreshTokenExpiration
    ) {
        return new LoginResponse(
                new UserDTO(user),
                accessToken,
                refreshToken,
                accessTokenExpiration,
                refreshTokenExpiration
        );
    }
}
