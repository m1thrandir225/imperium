package me.sebastijanzindl.authserver.controller;

import me.sebastijanzindl.authserver.dto.LoginUserDTO;
import me.sebastijanzindl.authserver.dto.RegisterUserDTO;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.responses.LoginResponse;
import me.sebastijanzindl.authserver.service.AuthenticationService;
import me.sebastijanzindl.authserver.service.JwtService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RequestMapping("/api/v1/auth")
@RestController
public class AuthController {
    private final JwtService jwtService;

    private final AuthenticationService authenticationService;

    public AuthController(JwtService jwtService, AuthenticationService authenticationService) {
        this.jwtService = jwtService;
        this.authenticationService = authenticationService;
    }

    @PostMapping("/register")
    public ResponseEntity<User> register(
            @RequestBody RegisterUserDTO registerUserDTO
    ) {
        User registeredUser = authenticationService.signup(registerUserDTO);
        return ResponseEntity.ok(registeredUser);
    }

    @PostMapping("/login")
    public ResponseEntity<LoginResponse> login(
            @RequestBody LoginUserDTO loginUserDTO
    ) {
        User user  = authenticationService.authenticate(loginUserDTO);

        String jwtToken = jwtService.generateToken(user);

        LoginResponse response = new LoginResponse();
        response.setToken(jwtToken);
        response.setExpiresIn(jwtService.getExpiration());

        return ResponseEntity.ok(response);
    }
}
