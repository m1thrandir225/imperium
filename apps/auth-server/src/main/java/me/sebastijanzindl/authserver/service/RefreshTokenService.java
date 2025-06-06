package me.sebastijanzindl.authserver.service;

import me.sebastijanzindl.authserver.model.RefreshToken;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.repository.RefreshTokenRepository;
import me.sebastijanzindl.authserver.repository.UserRepository;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.Optional;
import java.util.UUID;

@Service
public class RefreshTokenService {
    private final RefreshTokenRepository refreshTokenRepository;
    private final UserRepository userRepository;

    @Value("${security.jwt.refresh-token-expiration-key}")
    private Long refreshTokenDuration;

    public RefreshTokenService(RefreshTokenRepository refreshTokenRepository, UserRepository userRepository) {
        this.refreshTokenRepository = refreshTokenRepository;
        this.userRepository = userRepository;
    }

    public RefreshToken create(String userEmail) {
        User user = userRepository.findByEmail(userEmail).orElseThrow();

        RefreshToken refreshToken = RefreshToken.builder()
                .user(user)
                .token(UUID.randomUUID().toString())
                .expiresAt(Instant.now().plusMillis(refreshTokenDuration))
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
