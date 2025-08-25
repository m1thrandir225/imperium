package me.sebastijanzindl.authserver.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.Getter;
import lombok.Setter;

@Data
@AllArgsConstructor
public class UpdateUserPasswordDTO {
    @JsonProperty("password")
    @NotNull(message = "Password cannot be null")
    @NotBlank(message = "Password cannot be blank")
    private String password;

    @JsonProperty("new_password")
    @NotNull(message = "New password cannot be null")
    @NotBlank(message = "New password cannot be blank")
    private String newPassword;

}
