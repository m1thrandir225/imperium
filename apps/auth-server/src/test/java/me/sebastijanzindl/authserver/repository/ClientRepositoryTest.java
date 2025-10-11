package me.sebastijanzindl.authserver.repository;

import me.sebastijanzindl.authserver.model.Client;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.testsupport.PostgresTC;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;

import java.util.Optional;
import static org.junit.jupiter.api.Assertions.*;

@DataJpaTest
class ClientRepositoryIT extends PostgresTC {
    @Autowired
    private ClientRepository clientRepository;

    @Autowired
    private UserRepository userRepository;

    @Test
    void findByNameAndIpAddress_worksAgainstPostgres() {
        User u = new User();
        u.setName("Alice");
        u.setEmail("alice@example.com");
        u.setPassword("strongPassword");

        userRepository.save(u);

        Client client = new Client();

        client.setName("Random Client");
        client.setIpAddress("127.0.0.1");
        client.setOwner(u);

        clientRepository.save(client);

        Optional<Client> result = clientRepository.findByNameAndIpAddress(client.getName(), client.getIpAddress());

        assertTrue(result.isPresent());
        assertNotNull(result.get());
        assertNotNull(result.get().getOwner());
        assertNotNull(result.get().getIpAddress());
        assertNotNull(result.get().getId());

        //Client Verifications
        assertEquals(client.getIpAddress(), result.get().getIpAddress());
        assertEquals(client.getName(), result.get().getName());

        //Owner Verifications
        assertEquals(u.getName(), result.get().getOwner().getName());
        assertEquals(u.getEmail(), result.get().getOwner().getEmail());
    }

    @Test
    void findByNameAndIpAddress_returnsEmptyForNotFound() {
        assertFalse(clientRepository.findByNameAndIpAddress("Random Client", "127.0.0.1").isPresent());
    }

    @Test
    void findByNameAndIpAddress_returnsEmptyForNull() {
        assertFalse(clientRepository.findByNameAndIpAddress(null, null).isPresent());
    }
}
