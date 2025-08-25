package me.sebastijanzindl.authserver.service;

import jakarta.persistence.EntityNotFoundException;
import lombok.SneakyThrows;
import me.sebastijanzindl.authserver.dto.CreateClientDTO;
import me.sebastijanzindl.authserver.model.Client;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.repository.ClientRepository;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class ClientService {
    private final ClientRepository clientRepository;

    public ClientService(ClientRepository clientRepository) {
        this.clientRepository = clientRepository;
    }

    public Client create(CreateClientDTO input, User owner) {
        Client client = new Client();
        client.setName(input.getClientName());
        client.setIpAddress(input.getIpAddress());
        client.setOwner(owner);
        return clientRepository.save(client);
    }

    public Client getClient(UUID id){
        return this.findById(id);
    }

    @SneakyThrows
    public Client update(UUID id, CreateClientDTO input, User owner) {
        Client client = this.findById(id);

        if(!client.getOwner().equals(owner)) {
            throw new Exception("Not the same owner");
        }
        client.setName(input.getClientName());
        client.setIpAddress(input.getIpAddress());

        return clientRepository.save(client);
    }

    public Client delete(UUID id) {
        Client client = this.findById(id);

        clientRepository.delete(client);
        return client;
    }

    private Client findById(UUID id) {
        return this.clientRepository.findById(id).orElseThrow(() ->  new EntityNotFoundException("Client with id " + id + " not found"));
    }
}
