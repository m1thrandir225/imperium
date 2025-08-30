CREATE TABLE sessions
(
    id            UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    user_id       UUID         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    host_id       UUID         NOT NULL REFERENCES hosts (id) ON DELETE CASCADE,
    client_id     UUID         NOT NULL REFERENCES clients (id) ON DELETE CASCADE,
    status        VARCHAR(20)  NOT NULL DEFAULT 'PENDING',
    session_token VARCHAR(255) NOT NULL UNIQUE,
    webrtc_offer  TEXT,
    webrtc_answer TEXT,
    expires_at    TIMESTAMP    NOT NULL,
    created_at    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    started_at    TIMESTAMP,
    ended_at      TIMESTAMP,
    end_reason    TEXT,

    CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_sessions_host FOREIGN KEY (host_id) REFERENCES hosts (id),
    CONSTRAINT fk_sessions_client FOREIGN KEY (client_id) REFERENCES clients (id)
)

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_host_id ON sessions(host_id);
CREATE INDEX idx_sessions_client_id ON sessions(client_id);
CREATE INDEX idx_sessions_status ON sessions(status);
CREATE INDEX idx_sessions_token ON sessions(session_token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);