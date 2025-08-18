package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Date;

@Getter
@Setter
public class RefreshTokenResponse {
    @JsonProperty("access_token")
    private String accessToken;

    @JsonProperty("expires_at")
    private Date expiresAt;
}
