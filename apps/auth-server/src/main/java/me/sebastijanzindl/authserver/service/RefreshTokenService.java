package me.sebastijanzindl.authserver.service;

import me.sebastijanzindl.authserver.model.RefreshToken;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.TOKEN_TYPE;
import me.sebastijanzindl.authserver.repository.RefreshTokenRepository;
import me.sebastijanzindl.authserver.repository.UserRepository;
import me.sebastijanzindl.authserver.security.JwtUtils;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.util.Optional;

@Service
public class RefreshTokenService {
    private final RefreshTokenRepository refreshTokenRepository;
    private final JwtUtils jwtUtils;
    private final UserRepository userRepository;

    public RefreshTokenService(RefreshTokenRepository refreshTokenRepository, JwtUtils jwtUtils, UserRepository userRepository) {
        this.refreshTokenRepository = refreshTokenRepository;
        this.jwtUtils = jwtUtils;
        this.userRepository = userRepository;
    }

    public RefreshToken create(String userEmail) {
        User user = userRepository.findByEmail(userEmail).orElseThrow();

        RefreshToken refreshToken = RefreshToken.builder()
                .user(user)
                .token(jwtUtils.generateToken(user, TOKEN_TYPE.REFRESH))
                .expiresAt(Instant.now().plusMillis(jwtUtils.getTokenExpiration(TOKEN_TYPE.REFRESH)))
                .build();

        return refreshTokenRepository.save(refreshToken);
    }

    public Optional<RefreshToken> findByToken(String token) {
        return refreshTokenRepository.findByToken(token);
    }

    public RefreshToken verifyExpiration(RefreshToken token) {
        if(token.getExpiresAt().isBefore(Instant.now())) {
            refreshTokenRepository.delete(token);
            throw new RuntimeException(token.getToken() + " Refresh token expired. Please make a new signin request");
        }
        return token;
    }

}
