package me.sebastijanzindl.authserver.repository;

import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.testsupport.PostgresTC;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.context.annotation.Import;

import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;

@DataJpaTest
class UserRepositoryIT extends PostgresTC {

    @Autowired
    private UserRepository userRepository;

    @Test
    void saveAndFindByEmail_worksAgainstPostgres() {
        User u = new User();
        u.setEmail("alice@example.com");
        u.setPassword("strongPassword");
        u.setName("Alice Wonderland");

        userRepository.save(u);

        Optional<User> found = userRepository.findByEmail(u.getEmail());
        assertTrue(found.isPresent());
        assertNotNull(found.get().getId());
        assertNotNull(found.get().getPassword());
        assertNotNull(found.get().getCreatedAt());
        assertNotNull(found.get().getUpdatedAt());
        assertEquals(u.getEmail(), found.get().getEmail());
        assertEquals(u.getName(), found.get().getName());
    }

    @Test
    void findByEmail_returnsEmptyForMissingUser() {
        assertFalse(userRepository.findByEmail("bob@example.com").isPresent());
    }
}