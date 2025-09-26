package me.sebastijanzindl.authserver.security;

import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.TOKEN_TYPE;
import me.sebastijanzindl.authserver.testsupport.PostgresTC;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Profile;
import org.springframework.test.context.ActiveProfiles;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
@ActiveProfiles("test")
class JwtUtilsTest extends PostgresTC {

    @Autowired
    JwtUtils jwtUtils;

    @Test
    void generateAndValidateAccessToken() {
        User u = new User();
        u.setEmail("test@example.com");
        u.setPassword("x");
        u.setName("Test User");

        String token = jwtUtils.generateToken(u, TOKEN_TYPE.ACCESS);
        assertTrue(jwtUtils.isTokenValid(token, u, TOKEN_TYPE.ACCESS));
        assertEquals("test@example.com", jwtUtils.extractUsername(token, TOKEN_TYPE.ACCESS));
        assertNotNull(jwtUtils.extractExpiration(token, TOKEN_TYPE.ACCESS));
    }

    @Test
    void generateAndValidateRefreshToken() {
        User u = new User();
        u.setEmail("test@example.com");
        u.setPassword("x");
        u.setName("Test User");

        String token = jwtUtils.generateToken(u, TOKEN_TYPE.REFRESH);
        assertTrue(jwtUtils.isTokenValid(token, u, TOKEN_TYPE.REFRESH));
        assertEquals("test@example.com", jwtUtils.extractUsername(token, TOKEN_TYPE.REFRESH));
        assertNotNull(jwtUtils.extractExpiration(token, TOKEN_TYPE.REFRESH));
    }
}