package me.sebastijanzindl.authserver.controller;

import me.sebastijanzindl.authserver.dto.CreateHostDTO;
import me.sebastijanzindl.authserver.dto.UpdateHostDTO;
import me.sebastijanzindl.authserver.model.Host;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.HOST_STATUS;
import me.sebastijanzindl.authserver.responses.HostResponse;
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
    public ResponseEntity<HostResponse> createHost(
            @AuthenticationPrincipal User currentUser,
            @RequestBody CreateHostDTO createHostDTO
    ) {
        try {
            Host host = this.hostService.create(createHostDTO, currentUser);
            return ResponseEntity.ok(new HostResponse(host));
        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }

    @GetMapping("/{id}")
    public ResponseEntity<HostResponse> getHost(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id
            ) {
        try {
            Host host = this.hostService.getHost(id);
            return ResponseEntity.ok(new HostResponse(host));
        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }

    @GetMapping
    public ResponseEntity<List<HostResponse>> getUserHosts(
            @AuthenticationPrincipal User currentUser
    ) {
        List<Host> hosts = currentUser.getHosts();

        List<HostResponse> hostResponses = hosts.stream().map(HostResponse::new).toList();

        return ResponseEntity.ok(hostResponses);
    }

    @PutMapping("/{id}")
    public ResponseEntity<HostResponse> updateHost(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id,
            @RequestBody UpdateHostDTO updateHostDTO
            ) {
        try {
            Host updatedHost = this.hostService.update(id, updateHostDTO);
            return ResponseEntity.ok(new HostResponse(updatedHost));
        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }

    @PatchMapping("/{id}")
    public ResponseEntity<HostResponse> updateHostStatus(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id,
            @RequestBody HOST_STATUS hostStatus
    ) {
        try {
            Host updatedHost = this.hostService.updateStatus(id, hostStatus);
            return ResponseEntity.ok(new HostResponse(updatedHost));

        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<HostResponse> deleteHost(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id
    ) {
        try {
            Host deletedHost = this.hostService.delete(id);

            return ResponseEntity.ok(new HostResponse(deletedHost));
        } catch (Exception exception) {
            return ResponseEntity.badRequest().build();
        }
    }
}
