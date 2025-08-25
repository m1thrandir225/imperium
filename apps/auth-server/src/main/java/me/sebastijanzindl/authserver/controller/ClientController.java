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

        List<ClientResponse> clientResponses = clients.stream().map(ClientResponse::new).toList();
        return ResponseEntity.ok(clientResponses);
    }

    @GetMapping("/{id}")
    public ResponseEntity<ClientResponse> getClient(@PathVariable UUID id) {
        Client client = this.clientService.getClient(id);
        return ResponseEntity.ok(new ClientResponse(client));
    }

    @PostMapping
    public ResponseEntity<ClientResponse> createClient(
            @AuthenticationPrincipal User currentUser,
            @Valid @RequestBody CreateClientDTO input
    ) {
        Client client = this.clientService.create(input, currentUser);
        return ResponseEntity.ok(new ClientResponse(client));
    }

    @PutMapping("/{id}")
    public ResponseEntity<ClientResponse> updateClient(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id,
            @Valid @RequestBody CreateClientDTO input

    ) {
        Client client = this.clientService.update(id, input, currentUser);
        return ResponseEntity.ok(new ClientResponse(client));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteClient(
            @PathVariable UUID id
    ) {
        Client client = this.clientService.delete(id);
        return ResponseEntity.status(204).build();
    }
}
