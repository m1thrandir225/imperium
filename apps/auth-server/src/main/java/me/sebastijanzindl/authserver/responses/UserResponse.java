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


public record UserResponse(
        @JsonProperty("user")
        UserDTO user
) {
    public static UserResponse from(User user)
    {
        return new UserResponse(
                new UserDTO(user)
        );
    }

    public static List<UserResponse> fromList(List<User> users) {
        return users.stream()
                .map(UserResponse::from)
                .toList();
    }
}
