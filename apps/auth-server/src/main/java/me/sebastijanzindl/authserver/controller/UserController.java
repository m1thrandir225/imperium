package me.sebastijanzindl.authserver.controller;

import jakarta.validation.Valid;
import me.sebastijanzindl.authserver.dto.UpdateUserDTO;
import me.sebastijanzindl.authserver.dto.UpdateUserPasswordDTO;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.responses.UserResponse;
import me.sebastijanzindl.authserver.service.UserService;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

@RequestMapping("/api/v1/users")
@RestController
public class UserController {
    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping("/me")
    public ResponseEntity<UserResponse> getCurrentUser(
            @AuthenticationPrincipal User currentUser
    ) {
        return ResponseEntity.ok(new UserResponse(currentUser));
    }

    @PutMapping("/update")
    public ResponseEntity<UserResponse> updateUser(
            @AuthenticationPrincipal User user,
            @Valid @RequestBody UpdateUserDTO updateUserDTO
    ) {

        User updatedUser = this.userService.update(user.getId(), updateUserDTO);

        return ResponseEntity.ok(new UserResponse(updatedUser));
    }

    @PutMapping("/update-password")
    public ResponseEntity<UserResponse> updatePassword(
            @AuthenticationPrincipal User user,
            @Valid @RequestBody UpdateUserPasswordDTO dto
    ) {
        User updated = this.userService.updatePassword(user.getId(), dto.getPassword(), dto.getNewPassword());
        return ResponseEntity.ok(new UserResponse(updated));

    }

    @DeleteMapping
    public ResponseEntity<User> deleteUser(
            @AuthenticationPrincipal User user
    ) {
        User deletedUser = this.userService.delete(user.getId());
        return ResponseEntity.ok(deletedUser);
    }
}
