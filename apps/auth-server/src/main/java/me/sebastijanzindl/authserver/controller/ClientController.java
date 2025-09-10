package me.sebastijanzindl.authserver.controller;

import jakarta.validation.Valid;
import me.sebastijanzindl.authserver.dto.CreateClientDTO;
import me.sebastijanzindl.authserver.model.Client;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.responses.ClientResponse;
import me.sebastijanzindl.authserver.service.ClientService;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RequestMapping("/api/v1/clients")
@Controller
public class ClientController {
    private final ClientService clientService;

    public ClientController(ClientService clientService) {
        this.clientService = clientService;
    }

    @GetMapping("/me")
    public ResponseEntity<List<ClientResponse>> getUserClients(
            @AuthenticationPrincipal User currentUser
    ) {
        List<Client> clients = currentUser.getClients();

        List<ClientResponse> clientResponses = ClientResponse.fromList(clients);

        return ResponseEntity.ok(clientResponses);
    }

    @GetMapping("/{id}")
    public ResponseEntity<ClientResponse> getClient(@PathVariable UUID id) {
        Client client = this.clientService.getClient(id);

        ClientResponse response = ClientResponse.from(client);
        return ResponseEntity.ok(response);
    }

    @PostMapping
    public ResponseEntity<ClientResponse> createClient(
            @AuthenticationPrincipal User currentUser,
            @Valid @RequestBody CreateClientDTO input
    ) {
        Client client = this.clientService.create(input, currentUser);
        ClientResponse response = ClientResponse.from(client);
        return ResponseEntity.ok(response);
    }

    @PostMapping("/upsert")
    public ResponseEntity<ClientResponse> upsertClient(
            @AuthenticationPrincipal User currentUser,
            @Valid @RequestBody CreateClientDTO input
    ) {
        Client client = this.clientService.upsert(input, currentUser);
        ClientResponse response = ClientResponse.from(client);
        return ResponseEntity.ok(response);
    }

    @PutMapping("/{id}")
    public ResponseEntity<ClientResponse> updateClient(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id,
            @Valid @RequestBody CreateClientDTO input

    ) {
        Client client = this.clientService.update(id, input, currentUser);
        ClientResponse response = ClientResponse.from(client);

        return ResponseEntity.ok(response);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteClient(
            @PathVariable UUID id
    ) {
        Client client = this.clientService.delete(id);
        return ResponseEntity.status(204).build();
    }
}
