package me.sebastijanzindl.authserver.repository;

import me.sebastijanzindl.authserver.model.Host;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface HostRepository  extends JpaRepository<Host, UUID> {
    Optional<Host> findByName(String name);

    Optional<Host> findByNameAndIpAddressAndPort(String name, String ipAddress, Integer port);
}
