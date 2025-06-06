package me.sebastijanzindl.authserver.controller;

import me.sebastijanzindl.authserver.dto.CreateClientDTO;
import me.sebastijanzindl.authserver.model.Client;
import me.sebastijanzindl.authserver.model.User;
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

    @GetMapping
    public ResponseEntity<List<Client>> getUserClients(
            @AuthenticationPrincipal User currentUser
    ) {
        List<Client> clients = currentUser.getClients();
        return ResponseEntity.ok(clients);
    }

    @GetMapping("/{id}")
    public ResponseEntity<Client> getClient(@PathVariable UUID id) {
        try {
            Client client = this.clientService.getClient(id);
            return ResponseEntity.ok(client);
        } catch (Exception exception) {
            return ResponseEntity.notFound().build();
        }
    }

    @PostMapping
    public ResponseEntity<Client> createClient(
            @AuthenticationPrincipal User currentUser,
            @RequestBody CreateClientDTO input
    ) {
        try {
            Client client = this.clientService.create(input, currentUser);
            return ResponseEntity.ok(client);
        }  catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }

    @PutMapping
    public ResponseEntity<Client> updateClient(
            @AuthenticationPrincipal User currentUser,
            @RequestBody CreateClientDTO input,
            @PathVariable UUID id
    ) {
        try {
            Client client = this.clientService.update(id, input, currentUser);
            return ResponseEntity.ok(client);
        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }

    @DeleteMapping
    public ResponseEntity<Client> deleteClient(
            @PathVariable UUID id
    ) {
        try {
            Client client = this.clientService.delete(id);
            return ResponseEntity.ok(client);
        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }
}
