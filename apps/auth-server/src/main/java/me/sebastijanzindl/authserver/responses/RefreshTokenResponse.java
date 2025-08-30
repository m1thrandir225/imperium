package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import lombok.Setter;

import java.util.Date;


public record RefreshTokenResponse(
    @JsonProperty("access_token")
    String accessToken,

    @JsonProperty("expires_in")
    Date expiresIn
)
{
    public static RefreshTokenResponse from(String accessToken, Date expiresIn) {
        return new RefreshTokenResponse(accessToken, expiresIn);
    }
}
