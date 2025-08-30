package me.sebastijanzindl.authserver.service;

import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

@Service
public class SessionCleanupService {
    public final SessionService sessionService;

    public SessionCleanupService(SessionService sessionService) {
        this.sessionService = sessionService;
    }

    @Scheduled(fixedRate = 300000)
    public void cleanupExpiredSessions() {
        sessionService.cleanupExpiredSessions();
    }
}
