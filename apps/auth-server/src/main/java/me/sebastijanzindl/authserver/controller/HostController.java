package me.sebastijanzindl.authserver.controller;

import me.sebastijanzindl.authserver.dto.CreateHostDTO;
import me.sebastijanzindl.authserver.dto.UpdateHostDTO;
import me.sebastijanzindl.authserver.model.Host;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.HOST_STATUS;
import me.sebastijanzindl.authserver.service.AuthenticationService;
import me.sebastijanzindl.authserver.service.HostService;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RequestMapping("/api/v1/hosts")
@RestController
public class HostController {
    private final HostService hostService;
    private final AuthenticationService authenticationService;

    public HostController(HostService hostService, AuthenticationService authenticationService) {
        this.hostService = hostService;
        this.authenticationService = authenticationService;
    }

    @PostMapping
    public ResponseEntity<Host> createHost(
            @AuthenticationPrincipal User currentUser,
            @RequestBody CreateHostDTO createHostDTO
    ) {
        try {
            Host host = this.hostService.create(createHostDTO, currentUser);
            return ResponseEntity.ok(host);
        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }

    @GetMapping("/{id}")
    public ResponseEntity<Host> getHost(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id
            ) {
        try {
            Host host = this.hostService.getHost(id);
            return ResponseEntity.ok(host);
        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }

    @GetMapping
    public ResponseEntity<List<Host>> getUserHosts(
            @AuthenticationPrincipal User currentUser
    ) {
        List<Host> hosts = currentUser.getHosts();

        return ResponseEntity.ok(hosts);
    }

    @PutMapping("/{id}")
    public ResponseEntity<Host> updateHost(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id,
            @RequestBody UpdateHostDTO updateHostDTO
            ) {
        try {
            Host updatedHost = this.hostService.update(id, updateHostDTO);

            return ResponseEntity.ok(updatedHost);
        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }

    @PatchMapping("/{id}")
    public ResponseEntity<Host> updateHostStatus(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id,
            @RequestBody HOST_STATUS hostStatus
    ) {
        try {
            Host updatedHost = this.hostService.updateStatus(id, hostStatus);
            return ResponseEntity.ok(updatedHost);

        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Host> deleteHost(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id
    ) {
        try {
            Host deletedHost = this.hostService.delete(id);
            return ResponseEntity.ok(deletedHost);
        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }
}
