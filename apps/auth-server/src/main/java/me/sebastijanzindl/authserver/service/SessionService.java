package me.sebastijanzindl.authserver.service;

import jakarta.transaction.Transactional;
import me.sebastijanzindl.authserver.dto.CreateSessionDTO;
import me.sebastijanzindl.authserver.dto.EndSessionDTO;
import me.sebastijanzindl.authserver.dto.StartSessionDTO;
import me.sebastijanzindl.authserver.exceptions.SessionException;
import me.sebastijanzindl.authserver.model.Client;
import me.sebastijanzindl.authserver.model.Host;
import me.sebastijanzindl.authserver.model.Session;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.HOST_STATUS;
import me.sebastijanzindl.authserver.model.enums.SESSION_STATUS;
import me.sebastijanzindl.authserver.repository.ClientRepository;
import me.sebastijanzindl.authserver.repository.HostRepository;
import me.sebastijanzindl.authserver.repository.SessionRepository;
import me.sebastijanzindl.authserver.responses.SessionResponse;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.reactive.function.client.WebClientException;
import org.springframework.web.reactive.function.client.WebClientResponseException;

import java.net.http.HttpHeaders;
import java.time.LocalDateTime;
import java.util.*;

@Service
@Transactional
public class SessionService {
    private final SessionRepository sessionRepository;
    private final HostRepository hostRepository;
    private final ClientRepository clientRepository;
    private final WebClient webClient;

    public SessionService(
            SessionRepository sessionRepository,
            HostRepository hostRepository,
            ClientRepository clientRepository,
            WebClient.Builder webClientBuilder
    ) {
        this.sessionRepository = sessionRepository;
        this.hostRepository = hostRepository;
        this.clientRepository = clientRepository;
        this.webClient = webClientBuilder.build();
    }

    public Session createSession(CreateSessionDTO createSessionDTO, User user) {
        // Validate host exists and belongs to user
        Host host = hostRepository.findById(createSessionDTO.getHostId())
                .orElseThrow(() -> new SessionException("Host not found"));

        if (!host.getOwner().getId().equals(user.getId())) {
            throw new SessionException("Host does not belong to user");
        }

        // Validate client exists and belongs to user
        Client client = clientRepository.findById(createSessionDTO.getClientId())
                .orElseThrow(() -> new SessionException("Client not found"));

        if (!client.getOwner().getId().equals(user.getId())) {
            throw new SessionException("Client does not belong to user");
        }

        // Check if host is available
        if (host.getStatus() != HOST_STATUS.AVAILABLE) {
            throw new SessionException("Host is not available for sessions");
        }

        // Check if host already has an active session
        long activeSessions = sessionRepository.countByHostIdAndStatus(host.getId(), SESSION_STATUS.ACTIVE);
        if (activeSessions > 0) {
            throw new SessionException("Host already has an active session");
        }

        // Create session
        Session session = new Session();
        session.setUser(user);
        session.setHost(host);
        session.setClient(client);
        session.setStatus(SESSION_STATUS.PENDING);
        session.setSessionToken(generateSessionToken());
        session.setCreatedAt(LocalDateTime.now());
        session.setUpdatedAt(LocalDateTime.now());
        session.setProgramId(createSessionDTO.getProgramId());
        session.setExpiresAt(LocalDateTime.now().plusMinutes(60)); // 60-minute expiry

        return sessionRepository.save(session);
    }

    public Session startSession(UUID sessionId, StartSessionDTO startSessionDTO) {
        Session session = sessionRepository.findById(sessionId)
                .orElseThrow(() -> new SessionException("Session not found"));

        if (session.getStatus() != SESSION_STATUS.PENDING) {
            throw new SessionException("Session is not in pending status");
        }

        if (session.isExpired()) {
            session.cancel("Session expired before starting");
            sessionRepository.save(session);
            throw new SessionException("Session has expired");
        }

        session.setWebrtcOffer(startSessionDTO.getWebrtcOffer());
        try {
            String webrtcAnswer = startSessionOnHost(session.getHost(), session);

            session.setWebrtcAnswer(webrtcAnswer);
            session.start();
            // Update host status to INUSE
            Host host = session.getHost();
            host.setStatus(HOST_STATUS.INUSE);
            hostRepository.save(host);


            return sessionRepository.save(session);
        } catch(Exception e) {
            session.cancel("failed to start session: " + e.getMessage());
            sessionRepository.save(session);
            throw new SessionException("Failed to start session: " + e.getMessage());
        }
    }

    public Session endSession(UUID sessionId, EndSessionDTO endSessionDTO) {
        Session session = sessionRepository.findById(sessionId)
                .orElseThrow(() -> new SessionException("Session not found"));

        if (session.getStatus() != SESSION_STATUS.ACTIVE) {
            throw new SessionException("Session is not active");
        }

        if (endSessionDTO.getWebrtcAnswer() != null) {
            session.setWebrtcAnswer(endSessionDTO.getWebrtcAnswer());
        }

        session.end(endSessionDTO.getReason());

        try {
            endSessionOnHost(session.getHost(), session);
        } catch(Exception e) {
            System.err.println("Failed to end session on host: " + e.getMessage());
        }

        // Update host status back to AVAILABLE
        Host host = session.getHost();
        host.setStatus(HOST_STATUS.AVAILABLE);
        hostRepository.save(host);

        return sessionRepository.save(session);
    }

    public Session cancelSession(UUID sessionId, String reason) {
        Session session = sessionRepository.findById(sessionId)
                .orElseThrow(() -> new SessionException("Session not found"));

        if (session.getStatus() == SESSION_STATUS.ENDED || session.getStatus() == SESSION_STATUS.CANCELLED) {
            throw new SessionException("Session is already ended or cancelled");
        }

        session.cancel(reason);

        // If session was active, update host status
        if (session.getStatus() == SESSION_STATUS.ACTIVE) {
            try {
                endSessionOnHost(session.getHost(), session);
            } catch(Exception e) {
                System.err.println("Failed to end session on host: " + e.getMessage());
            }
            Host host = session.getHost();
            host.setStatus(HOST_STATUS.AVAILABLE);
            hostRepository.save(host);
        }

        return sessionRepository.save(session);
    }

    public Session getSession(UUID sessionId) {
        return sessionRepository.findById(sessionId)
                .orElseThrow(() -> new SessionException("Session not found"));
    }

    public Session getSessionByToken(String sessionToken) {
        return sessionRepository.findBySessionToken(sessionToken)
                .orElseThrow(() -> new SessionException("Invalid session token"));
    }

    public List<Session> getUserSessions(User user) {
        return sessionRepository.findByUserOrderByCreatedAtDesc(user);
    }

    public List<Session> getHostSessions(UUID hostId) {
        return sessionRepository.findByHostIdOrderByCreatedAtDesc(hostId);
    }

    public List<Session> getClientSessions(UUID clientId) {
        return sessionRepository.findByClientIdOrderByCreatedAtDesc(clientId);
    }

    public boolean validateSessionToken(String sessionToken, UUID hostId) {
        Optional<Session> sessionOpt = sessionRepository.findBySessionToken(sessionToken);
        if (sessionOpt.isEmpty()) {
            return false;
        }

        Session session = sessionOpt.get();
        return session.getHost().getId().equals(hostId) &&
                session.getStatus() == SESSION_STATUS.ACTIVE &&
                !session.isExpired();
    }

    private String generateSessionToken() {
        return "session_" + UUID.randomUUID().toString().replace("-", "");
    }

    public void cleanupExpiredSessions() {
        List<Session> expiredSessions = sessionRepository.findExpiredSessions(
                LocalDateTime.now(),
                List.of(SESSION_STATUS.PENDING, SESSION_STATUS.ACTIVE)
        );

        for (Session session : expiredSessions) {
            session.cancel("Session expired");
            sessionRepository.save(session);

            if (session.getStatus() == SESSION_STATUS.ACTIVE) {
                Host host = session.getHost();
                host.setStatus(HOST_STATUS.AVAILABLE);
                hostRepository.save(host);
            }
        }
    }

    public String startSessionOnHost(Host host, Session session) {
        String hostUrl = "http://" + host.getIpAddress() + ":" + host.getPort();
        String endpoint = hostUrl + "/api/session/start";

        SessionResponse sessionResponse = SessionResponse.from(session);

        try {
            Map<String, Object> response = this.webClient
                    .post()
                    .uri(endpoint)
                    .bodyValue(sessionResponse)
                    .retrieve()
                    .bodyToMono(Map.class)
                    .block();
            if (response != null && response.containsKey("webrtc_answer")) {
                return (String) response.get("webrtc_answer");
            } else {
                throw new RuntimeException("Failed to start session on host");
            }
        } catch(WebClientResponseException e) {
            throw new RuntimeException("Host returned error: " + e.getStatusCode() + " - " + e.getResponseBodyAsString());
        } catch (Exception e) {
            throw new RuntimeException("failed to communicate with host: " + e.getMessage(), e);
        }
    }

    public void endSessionOnHost(Host host, Session session) {
        String hostUrl = "http://" + host.getIpAddress() + ":" + host.getPort();
        String endpoint = hostUrl + "/api/session/end";

        SessionResponse sessionResponse = SessionResponse.from(session);

        try {
            webClient
                    .post()
                    .uri(endpoint)
                    .bodyValue(sessionResponse)
                    .retrieve()
                    .bodyToMono(Map.class)
                    .block();
        } catch (WebClientResponseException e) {
            throw new RuntimeException("Host returned error: " + e.getStatusCode() + " - " + e.getResponseBodyAsString());
        } catch (Exception e) {
            throw new RuntimeException("failed to communicate with host: " + e.getMessage(), e);
        }
    }
}
