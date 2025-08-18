package me.sebastijanzindl.authserver.security;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.io.Decoders;
import io.jsonwebtoken.security.Keys;
import me.sebastijanzindl.authserver.model.enums.TOKEN_TYPE;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;

import javax.crypto.SecretKey;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.function.Function;

@Service
public class JwtUtils {
    @Value("${security.jwt.access-token-security-key}")
    private String accessTokenSecurityKey;

    @Value("${security.jwt.access-token-expiration-time}")
    private Long accessTokenExpiration;

    @Value("${security.jwt.refresh-token-security-key}")
    private String refreshTokenSecurityKey;

    @Value("${security.jwt.refresh-token-expiration-time}")
    private Long refreshTokenExpiration;


    public Long getTokenExpiration(TOKEN_TYPE type) throws IllegalArgumentException{
        return switch (type) {
            case ACCESS -> accessTokenExpiration;
            case REFRESH -> refreshTokenExpiration;
            default -> throw new IllegalArgumentException("Invalid token type");
        };
    }

    public String generateToken(UserDetails userDetails, TOKEN_TYPE type) {
        return generateToken(new HashMap<>(), userDetails, type);
    }

    public String generateToken(Map<String, Object> extraClaims, UserDetails userDetails, TOKEN_TYPE type) {
        return Jwts.builder()
                .claims(extraClaims)
                .subject(userDetails.getUsername())
                .issuedAt(new Date(System.currentTimeMillis()))
                .expiration(new Date(System.currentTimeMillis() + getTokenExpiration(type)))
                .signWith(getSignInKey(type))
                .compact();
    }

    public boolean isTokenValid(String token, UserDetails userDetails, TOKEN_TYPE type) {
        final String username = extractUsername(token, type);
        return (username.equals(userDetails.getUsername())) && !isTokenExpired(token, type);
    }

    public String extractUsername(String token, TOKEN_TYPE type) {
        return extractClaim(token, Claims::getSubject, type);
    }

    public <T> T extractClaim(String token, Function<Claims, T> claimsResolver, TOKEN_TYPE type) {
        final Claims claims = extractAllClaims(token, type);
        return claimsResolver.apply(claims);
    }

    private boolean isTokenExpired(String token, TOKEN_TYPE type) {
        return extractExpiration(token, type).before(new Date());
    }

    private Date extractExpiration(String token,  TOKEN_TYPE type) {
        return extractClaim(token, Claims::getExpiration, type);
    }

    private Claims extractAllClaims(String token, TOKEN_TYPE type) {
        return Jwts.parser()
                .verifyWith(getSignInKey(type))
                .build()
                .parseSignedClaims(token)
                .getPayload();
    }

    private SecretKey getSignInKey(TOKEN_TYPE type) {
        byte[] keyBytes;
        switch (type) {
            case ACCESS -> keyBytes = Decoders.BASE64.decode(accessTokenSecurityKey);
            case REFRESH -> keyBytes = Decoders.BASE64.decode(refreshTokenSecurityKey);
            default -> throw new IllegalStateException("Unexpected value: " + type);
        }
        return Keys.hmacShaKeyFor(keyBytes);
    }
}
