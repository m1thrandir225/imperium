package me.sebastijanzindl.authserver.repository;

import me.sebastijanzindl.authserver.testsupport.PostgresTC;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;

@DataJpaTest
class SessionRepositoryIT extends PostgresTC {
    @Autowired
    private SessionRepository sessionRepository;

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private ClientRepository clientRepository;

    @Autowired
    private HostRepository hostRepository;

    @BeforeEach
    void init() {
    }

    @Test
    void findByUserOrderByCreatedAtDesc_worksAgainstPostgres() {
    }

    @Test
    void findBySessionToken_worksAgainstPostgres() {}

    @Test
    void findActiveSessionsByStatus_worksAgainstPostgres() {}

    @Test
    void findByHostIdAndStatus_worksAgainstPostgres() {}

    @Test
    void findByHostIdOrderByCreatedAtDesc_worksAgainstPostgres() {}

    @Test
    void findyByClientIdOrderByCreatedAtDesc_worksAgainstPostgres() {}

    @Test
    void countByHostIdAndStatus_worksAgainstPostgres() {}

    @Test
    void findExpiredSessions_worksAgainstPostgres() {}
}
