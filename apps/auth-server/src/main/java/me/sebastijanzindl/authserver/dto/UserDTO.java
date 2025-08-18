package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import me.sebastijanzindl.authserver.model.User;

import java.util.Date;
import java.util.UUID;

@Getter
public class UserDTO {
    @JsonProperty("id")
    private UUID id;

    @JsonProperty("email")
    private String email;

    @JsonProperty("name")
    private String name;

    @JsonProperty("created_at")
    private Date createdAt;

    @JsonProperty("updated_at")
    private Date updatedAt;

    public UserDTO(User user) {
        this.id = user.getId();
        this.email = user.getEmail();
        this.name = user.getName();
        this.createdAt = user.getCreatedAt();
        this.updatedAt = user.getUpdatedAt();
    }
}
