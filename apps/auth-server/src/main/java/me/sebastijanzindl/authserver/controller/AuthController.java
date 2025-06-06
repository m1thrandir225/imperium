package me.sebastijanzindl.authserver.controller;

import me.sebastijanzindl.authserver.dto.LoginUserDTO;
import me.sebastijanzindl.authserver.dto.RegisterUserDTO;
import me.sebastijanzindl.authserver.model.RefreshToken;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.responses.LoginResponse;
import me.sebastijanzindl.authserver.dto.RefreshTokenDTO;
import me.sebastijanzindl.authserver.responses.RefreshTokenResponse;
import me.sebastijanzindl.authserver.responses.UserResponse;
import me.sebastijanzindl.authserver.service.AuthenticationService;
import me.sebastijanzindl.authserver.service.JwtService;
import me.sebastijanzindl.authserver.service.RefreshTokenService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RequestMapping("/api/v1/auth")
@RestController
public class AuthController {
    private final JwtService jwtService;
    private final RefreshTokenService refreshTokenService;
    private final AuthenticationService authenticationService;

    public AuthController(JwtService jwtService, RefreshTokenService refreshTokenService, AuthenticationService authenticationService) {
        this.jwtService = jwtService;
        this.refreshTokenService = refreshTokenService;
        this.authenticationService = authenticationService;
    }

    @PostMapping("/register")
    public ResponseEntity<UserResponse> register(
            @RequestBody RegisterUserDTO registerUserDTO
    ) {
        User registeredUser = authenticationService.signup(registerUserDTO);
        return ResponseEntity.ok(new UserResponse(registeredUser));
    }

    @PostMapping("/login")
    public ResponseEntity<LoginResponse> login(
            @RequestBody LoginUserDTO loginUserDTO
    ) {
        User user  = authenticationService.authenticate(loginUserDTO);

        String jwtToken = jwtService.generateToken(user);

        RefreshToken refreshToken = refreshTokenService.create(user.getEmail());

        LoginResponse response = new LoginResponse();
        response.setToken(jwtToken);
        response.setExpiresIn(jwtService.getExpiration());
        response.setRefreshToken(refreshToken.getToken());

        return ResponseEntity.ok(response);
    }

    @PostMapping("/refresh")
    public ResponseEntity<RefreshTokenResponse> refreshToken(
            @RequestBody RefreshTokenDTO refreshTokenDTO
    ) {
        return refreshTokenService.findByToken(refreshTokenDTO.getToken())
                .map(refreshTokenService::verifyExpiration)
                .map(RefreshToken::getUser)
                .map(user -> {
                    String jwtToken = jwtService.generateToken(user);
                    RefreshTokenResponse response = new RefreshTokenResponse(jwtToken);
                    return ResponseEntity.ok(response);
                }).orElseThrow(() -> new RuntimeException("Refresh token not found"));
    }
}
