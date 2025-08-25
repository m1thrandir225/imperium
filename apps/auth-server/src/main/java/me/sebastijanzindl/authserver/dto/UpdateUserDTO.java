package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.Getter;

@Data
@AllArgsConstructor
public class UpdateUserDTO {
    @JsonProperty("email")
    @NotNull(message = "Email cannot be null")
    @NotBlank(message = "Email cannot be blank")
    private String email;

    @JsonProperty("first_name")
    @NotNull(message = "First name cannot be null")
    private String firstName;

    @JsonProperty("last_name")
    @NotNull(message = "Last name cannot be null")
    private String lastName;
}
