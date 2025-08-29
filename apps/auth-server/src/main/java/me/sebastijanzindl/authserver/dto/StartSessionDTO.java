package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import lombok.Data;

@Data
public class StartSessionDTO {
    @JsonProperty("webrtc_offer")
    @NotBlank(message = "WebRTC offer cannot be blank")
    private String webrtcOffer;
}
