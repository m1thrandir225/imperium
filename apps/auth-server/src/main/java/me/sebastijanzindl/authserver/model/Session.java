package me.sebastijanzindl.authserver.model;

import jakarta.persistence.*;
import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import me.sebastijanzindl.authserver.model.enums.SESSION_STATUS;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.LocalDateTime;
import java.util.UUID;

@Entity
@Getter
@Setter
@EqualsAndHashCode
@NoArgsConstructor
@Table(name = "sessions")
public class Session {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    @Column(columnDefinition = "uuid", updatable = false, nullable = false)
    private UUID id;

    @ManyToOne
    @JoinColumn(name = "user_id", nullable = false)
    private User user;

    @ManyToOne
    @JoinColumn(name = "host_id", nullable = false)
    private Host host;

    @ManyToOne
    @JoinColumn(name = "client_id", nullable = false)
    private Client client;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private SESSION_STATUS status = SESSION_STATUS.PENDING;

    @Column(nullable = false, unique = true, name = "session_token")
    private String sessionToken;

    @Column(name = "program_id")
    private String programId;

    @Column(columnDefinition = "TEXT", name = "webrtc_offer")
    private String webrtcOffer;

    @Column(columnDefinition = "TEXT", name = "webrtc_answer")
    private String webrtcAnswer;

    @Column(nullable = false, name = "expires_at")
    private LocalDateTime expiresAt;

    @CreationTimestamp
    @Column(updatable = false, name = "created_at")
    private LocalDateTime createdAt;

    @UpdateTimestamp
    @Column(name = "updated_at")
    private LocalDateTime updatedAt;

    @Column(name = "started_at")
    private LocalDateTime startedAt;

    @Column(name = "ended_at")
    private LocalDateTime endedAt;

    @Column(name = "end_reason", columnDefinition = "TEXT")
    private String endReason;

    public boolean isActive() {
        return this.status == SESSION_STATUS.ACTIVE;
    }

    public boolean isExpired() {
        return this.expiresAt.isBefore(LocalDateTime.now());
    }

    public void start() {
        this.status = SESSION_STATUS.ACTIVE;
        this.startedAt = LocalDateTime.now();
    }

    public void end(String reason) {
        this.status = SESSION_STATUS.ENDED;
        this.endedAt = LocalDateTime.now();
        this.endReason = reason;
    }

    public void cancel(String reason) {
        this.status = SESSION_STATUS.CANCELLED;
        this.endedAt = LocalDateTime.now();
        this.endReason = reason;
    }
}
