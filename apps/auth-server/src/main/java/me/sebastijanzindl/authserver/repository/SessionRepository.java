package me.sebastijanzindl.authserver.repository;

import me.sebastijanzindl.authserver.model.Session;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.SESSION_STATUS;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

public interface SessionRepository extends JpaRepository<Session, UUID> {
    List<Session> findByUserOrderByCreatedAtDesc(User user);

    List<Session> findByHostIdOrderByCreatedAtDesc(UUID hostId);

    List<Session> findByClientIdOrderByCreatedAtDesc(UUID clientId);

    Optional<Session> findBySessionToken(String sessionToken);

    @Query("SELECT s FROM Session s WHERE s.status = :status AND s.expiresAt > :now")
    List<Session> findActiveSessionsByStatus(@Param("status") SESSION_STATUS status, @Param("now") LocalDateTime now);

    @Query("SELECT s FROM Session s WHERE s.host.id = :hostId AND s.status = :status")
    List<Session> findByHostIdAndStatus(@Param("hostId") UUID hostId, @Param("status") SESSION_STATUS status);

    @Query("SELECT COUNT(s) FROM Session s WHERE s.host.id = :hostId AND s.status = :status")
    long countByHostIdAndStatus(@Param("hostId") UUID hostId, @Param("status") SESSION_STATUS status);

    @Query("SELECT s FROM Session s WHERE s.expiresAt < :now AND s.status IN (:statuses)")
    List<Session> findExpiredSessions(@Param("now") LocalDateTime now, @Param("statuses") List<SESSION_STATUS> statuses);
}
