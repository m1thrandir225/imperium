package me.sebastijanzindl.authserver.controller;

import jakarta.validation.Valid;
import me.sebastijanzindl.authserver.dto.CreateSessionDTO;
import me.sebastijanzindl.authserver.dto.EndSessionDTO;
import me.sebastijanzindl.authserver.dto.StartSessionDTO;
import me.sebastijanzindl.authserver.model.Session;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.responses.SessionResponse;
import me.sebastijanzindl.authserver.responses.SimpleSessionResponse;
import me.sebastijanzindl.authserver.service.SessionService;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

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
        SessionResponse response = new SessionResponse(session);
        return ResponseEntity.ok(response);
    }

    @PostMapping("/{sessionId}/start")
    public ResponseEntity<SessionResponse> startSession(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID sessionId,
            @Valid @RequestBody StartSessionDTO startSessionDTO
    ) {
        Session session = sessionService.startSession(sessionId, startSessionDTO);
        SessionResponse response = new SessionResponse(session);
        return ResponseEntity.ok(response);
    }

    @PostMapping("/{sessionId}/end")
    public ResponseEntity<SessionResponse> endSession(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID sessionId,
            @RequestBody EndSessionDTO endSessionDTO
    ) {
        Session session = sessionService.endSession(sessionId, endSessionDTO);
        SessionResponse response = new SessionResponse(session);
        return ResponseEntity.ok(response);
    }

    @PostMapping("/{sessionId}/cancel")
    public ResponseEntity<SessionResponse> cancelSession(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID sessionId,
            @RequestParam(defaultValue = "Cancelled by user") String reason
    ) {
        Session session = sessionService.cancelSession(sessionId, reason);
        SessionResponse response = new SessionResponse(session);
        return ResponseEntity.ok(response);
    }

    @GetMapping("/{sessionId}")
    public ResponseEntity<SessionResponse> getSession(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID sessionId
    ) {
        Session session = sessionService.getSession(sessionId);
        SessionResponse response = new SessionResponse(session);
        return ResponseEntity.ok(response);
    }

    @GetMapping
    public ResponseEntity<List<SimpleSessionResponse>> getUserSessions(
            @AuthenticationPrincipal User currentUser
    ) {
        List<Session> sessions = sessionService.getUserSessions(currentUser);
        List<SimpleSessionResponse> responses = sessions.stream()
                .map(SimpleSessionResponse::new)
                .toList();
        return ResponseEntity.ok(responses);
    }

    @GetMapping("/host/{hostId}")
    public ResponseEntity<List<SimpleSessionResponse>> getHostSessions(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID hostId
    ) {
        List<Session> sessions = sessionService.getHostSessions(hostId);
        List<SimpleSessionResponse> responses = sessions.stream()
                .map(SimpleSessionResponse::new)
                .toList();
        return ResponseEntity.ok(responses);
    }

    @GetMapping("/client/{clientId}")
    public ResponseEntity<List<SimpleSessionResponse>> getClientSessions(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID clientId
    ) {
        List<Session> sessions = sessionService.getClientSessions(clientId);
        List<SimpleSessionResponse> responses = sessions.stream()
                .map(SimpleSessionResponse::new)
                .toList();
        return ResponseEntity.ok(responses);
    }

    @PostMapping("/validate")
    public ResponseEntity<Boolean> validateSessionToken(
            @RequestParam String sessionToken,
            @RequestParam UUID hostId
    ) {
        boolean isValid = sessionService.validateSessionToken(sessionToken, hostId);
        return ResponseEntity.ok(isValid);
    }
}
