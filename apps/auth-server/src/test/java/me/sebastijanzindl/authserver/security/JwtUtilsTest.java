package me.sebastijanzindl.authserver.security;

import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.TOKEN_TYPE;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.runner.ApplicationContextRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.*;

class JwtUtilsTest {
    private final ApplicationContextRunner contextRunner = new ApplicationContextRunner()
            .withPropertyValues(
                    "security.jwt.access-token-security-key=UW5HOEdUTXBGY1FMeGdpNVRUNVdhM2EzMjQ3ODJLWm8=",
                    "security.jwt.refresh-token-security-key=UW5HOEdUTXBGY1FMeGdpNVRUNVdhM2EzMjQ3ODJLWm8=",
                    "security.jwt.access-token-expiration-time=600000",
                    "security.jwt.refresh-token-expiration-time=1209600000"
            )
            .withUserConfiguration(TestConfig.class);

    @Configuration
    static class TestConfig {
        @Bean
        JwtUtils jwtUtils() {
            return new JwtUtils();
        }
    }

    @Test
    void generateAndValidateAccessToken() {
        contextRunner.run(ctx -> {
            JwtUtils jwtUtils = ctx.getBean(JwtUtils.class);

            User u = new User();
            u.setEmail("test@example.com");
            u.setPassword("x");
            u.setName("Test User");

            String token = jwtUtils.generateToken(u, TOKEN_TYPE.ACCESS);

            assertThat(jwtUtils.isTokenValid(token, u, TOKEN_TYPE.ACCESS)).isTrue();
            assertThat(jwtUtils.extractUsername(token, TOKEN_TYPE.ACCESS)).isEqualTo("test@example.com");
            assertThat(jwtUtils.extractExpiration(token, TOKEN_TYPE.ACCESS)).isNotNull();
        });
    }

    @Test
    void generateAndValidateRefreshToken() {
        contextRunner.run(ctx -> {
            JwtUtils jwtUtils = ctx.getBean(JwtUtils.class);

            User u = new User();
            u.setEmail("test@example.com");
            u.setPassword("x");
            u.setName("Test User");

            String token = jwtUtils.generateToken(u, TOKEN_TYPE.REFRESH);
            assertTrue(jwtUtils.isTokenValid(token, u, TOKEN_TYPE.REFRESH));
            assertEquals("test@example.com", jwtUtils.extractUsername(token, TOKEN_TYPE.REFRESH));
            assertNotNull(jwtUtils.extractExpiration(token, TOKEN_TYPE.REFRESH));
        });

    }
}