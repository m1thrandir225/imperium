package me.sebastijanzindl.authserver.repository;

import me.sebastijanzindl.authserver.model.Host;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.HOST_STATUS;
import me.sebastijanzindl.authserver.responses.HostResponse;
import me.sebastijanzindl.authserver.testsupport.PostgresTC;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;

import java.util.Optional;

@DataJpaTest
class HostRepositoryIT extends PostgresTC {
    @Autowired
    private HostRepository hostRepository;

    @Autowired
    private UserRepository userRepository;

    private User createRandomUser() {
        User user = new User();

        user.setName("Alice");
        user.setPassword("strongPassword");
        user.setEmail("alice@example.com");

        userRepository.save(user);

        return user;
    }

    @Test
    void findByName_worksAgainstPostgres() {
        User user = createRandomUser();
        Host host = new Host();

        host.setName("Alice Host");
        host.setPort(8001);
        host.setStatus(HOST_STATUS.AVAILABLE);
        host.setIpAddress("127.0.0.1");
        host.setOwner(user);
        hostRepository.save(host);

        Host result = hostRepository.findByName(host.getName());

        assert
    }

    @Test
    void findByName_returnsEmptyForNotFound() {

    }

    @Test
    void findByNameAndIpAddressAndPort_worksAgainstPostgres() {
        User user = createRandomUser();
        Host host = new Host();

        host.setName("Alice Host");
        host.setPort(8001);
        host.setStatus(HOST_STATUS.AVAILABLE);
        host.setIpAddress("127.0.0.1");
        host.setOwner(user);
        hostRepository.save(host);
    }

    @Test
    void findByNameAndIpAddressAndPort_returnsEmptyForNotFound() {}
}
