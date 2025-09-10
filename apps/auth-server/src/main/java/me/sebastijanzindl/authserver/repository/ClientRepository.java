package me.sebastijanzindl.authserver.repository;

import jakarta.transaction.Transactional;
import me.sebastijanzindl.authserver.model.Client;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface ClientRepository extends JpaRepository<Client, UUID> {
    @Transactional
    default Client updateOrInsert(Client client) {
        return this.save(client);
    }

    Optional<Client> findByNameAndIpAddress(String name, String ipAddress);
}
