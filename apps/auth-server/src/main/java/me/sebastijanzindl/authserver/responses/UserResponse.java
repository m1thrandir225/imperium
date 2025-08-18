package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;
import me.sebastijanzindl.authserver.dto.UserDTO;
import me.sebastijanzindl.authserver.model.Client;
import me.sebastijanzindl.authserver.model.Host;
import me.sebastijanzindl.authserver.model.User;

import java.util.Date;
import java.util.List;
import java.util.UUID;

@Getter
@Setter
public class UserResponse {
    @JsonProperty("user")
    public UserDTO user;

    public UserResponse(User user) {
        this.user = new UserDTO(user);
    }
}
