package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class EndSessionDTO {
    @JsonProperty("reason")
    private String reason;

    @JsonProperty("webrtc_answer")
    private String webrtcAnswer;
}
