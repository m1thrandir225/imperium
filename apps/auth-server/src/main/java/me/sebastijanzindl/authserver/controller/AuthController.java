package me.sebastijanzindl.authserver.controller;

import me.sebastijanzindl.authserver.dto.LoginUserDTO;
import me.sebastijanzindl.authserver.dto.RegisterUserDTO;
import me.sebastijanzindl.authserver.model.RefreshToken;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.TOKEN_TYPE;
import me.sebastijanzindl.authserver.responses.LoginResponse;
import me.sebastijanzindl.authserver.dto.RefreshTokenDTO;
import me.sebastijanzindl.authserver.responses.RefreshTokenResponse;
import me.sebastijanzindl.authserver.responses.UserResponse;
import me.sebastijanzindl.authserver.service.AuthenticationService;
import me.sebastijanzindl.authserver.security.JwtUtils;
import me.sebastijanzindl.authserver.service.RefreshTokenService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RequestMapping("/api/v1/auth")
@RestController
public class AuthController {
    private final JwtUtils jwtUtils;
    private final RefreshTokenService refreshTokenService;
    private final AuthenticationService authenticationService;

    public AuthController(
            JwtUtils jwtUtils,
            RefreshTokenService refreshTokenService,
            AuthenticationService authenticationService
    ) {
        this.jwtUtils = jwtUtils;
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

        String jwtToken = jwtUtils.generateToken(user, TOKEN_TYPE.ACCESS);
        String refreshToken = jwtUtils.generateToken(user, TOKEN_TYPE.REFRESH);


        LoginResponse response = new LoginResponse();
        response.setToken(jwtToken);
        response.setRefreshToken(refreshToken);
        response.setExpiresIn(jwtUtils.getTokenExpiration(TOKEN_TYPE.ACCESS));

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
                    String jwtToken = jwtUtils.generateToken(user, TOKEN_TYPE.ACCESS);
                    RefreshTokenResponse response = new RefreshTokenResponse(jwtToken);
                    return ResponseEntity.ok(response);
                }).orElseThrow(() -> new RuntimeException("Refresh token not found"));
    }
}
