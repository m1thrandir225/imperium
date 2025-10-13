package me.sebastijanzindl.authserver.repository;

import me.sebastijanzindl.authserver.model.RefreshToken;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.testsupport.PostgresTC;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;

import java.time.Instant;
import java.util.Optional;

@DataJpaTest
class RefreshTokenRepositoryIT extends PostgresTC {
    @Autowired
    private RefreshTokenRepository refreshTokenRepository;

    @Autowired
    private UserRepository userRepository;


    private User createUser() {
        User user = new User();
        user.setName("Alice");
        user.setPassword("strongPassword");
        user.setEmail("alice@example.com");

        userRepository.save(user);
        return user;
    }

    private RefreshToken createRefreshToken() {
        User user = createUser();
        RefreshToken refreshToken = RefreshToken.builder()
                .user(user)
                .token("example_token")
                .expiresAt(Instant.now().plusMillis(60000))
                .build();
        refreshTokenRepository.save(refreshToken);

        return refreshToken;
    }

    @Test
    void findByToken_worksAgainstPostgres() {
        RefreshToken refreshToken = createRefreshToken();

        Optional<RefreshToken> result = refreshTokenRepository.findByToken(refreshToken.getToken());

        Assertions.assertTrue(result.isPresent());
        Assertions.assertEquals(refreshToken.getToken(), result.get().getToken());
        Assertions.assertEquals(refreshToken.getExpiresAt(), result.get().getExpiresAt());
    }

    @Test
    void findByTokenNotFound_worksAgainstPostgres() {
        Assertions.assertFalse(refreshTokenRepository.findByToken("token").isPresent());
    }
}
