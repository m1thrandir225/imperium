package me.sebastijanzindl.authserver.controller;

import jakarta.validation.Valid;
import me.sebastijanzindl.authserver.dto.CreateSessionDTO;
import me.sebastijanzindl.authserver.model.Session;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.responses.SessionResponse;
import me.sebastijanzindl.authserver.service.SessionService;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RequestMapping("/api/v1/sessions")
@RestController
public class SessionController {
    private final SessionService sessionService;

    public SessionController(SessionService sessionService) {
        this.sessionService = sessionService;
    }

    @PostMapping
    public ResponseEntity<SessionResponse> createSession(
            @AuthenticationPrincipal User currentUser,
            @Valid @RequestBody CreateSessionDTO createSessionDTO
    ) {
        Session session = sessionService.createSession(createSessionDTO, currentUser);
        SessionResponse response = new SessionResponse(
                session.getId(),
                session.getHost().getId(),
                session.getClient().getId()...
        );
        return ResponseEntity.ok(response);
    }
}
