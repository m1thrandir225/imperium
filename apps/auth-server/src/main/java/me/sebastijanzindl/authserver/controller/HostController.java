package me.sebastijanzindl.authserver.controller;

import jakarta.validation.Valid;
import me.sebastijanzindl.authserver.dto.CreateHostDTO;
import me.sebastijanzindl.authserver.dto.UpdateHostDTO;
import me.sebastijanzindl.authserver.model.Host;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.HOST_STATUS;
import me.sebastijanzindl.authserver.responses.HostResponse;
import me.sebastijanzindl.authserver.responses.SimpleHostResponse;
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
            @Valid @RequestBody CreateHostDTO createHostDTO
    ) {
        Host host = this.hostService.create(createHostDTO, currentUser);
        HostResponse response = HostResponse.from(host);
        return ResponseEntity.ok(response);
    }

    @GetMapping("/{id}")
    public ResponseEntity<HostResponse> getHost(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id
    ){
        Host host = this.hostService.getHost(id);
        HostResponse response = HostResponse.from(host);
        return ResponseEntity.ok(response);
    }

    @GetMapping
    public ResponseEntity<List<SimpleHostResponse>> getUserHosts(
            @AuthenticationPrincipal User currentUser
    ) {
        List<Host> hosts = currentUser.getHosts();

        List<SimpleHostResponse> hostResponses = SimpleHostResponse.fromList(hosts);

        return ResponseEntity.ok(hostResponses);
    }

    @PutMapping("/{id}")
    public ResponseEntity<HostResponse> updateHost(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id,
            @Valid @RequestBody UpdateHostDTO updateHostDTO
    ) {
        Host updatedHost = this.hostService.update(id, updateHostDTO);
        HostResponse response = HostResponse.from(updatedHost);

        return ResponseEntity.ok(response);
    }

    @PatchMapping("/{id}/status")
    public ResponseEntity<HostResponse> updateHostStatus(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id,
            @RequestBody HOST_STATUS hostStatus
    ) {
        Host updatedHost = this.hostService.updateStatus(id, hostStatus);
        HostResponse response = HostResponse.from(updatedHost);
        return ResponseEntity.ok(response);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteHost(
            @AuthenticationPrincipal User currentUser,
            @PathVariable UUID id
    ) {
        Host deletedHost = this.hostService.delete(id);

        return ResponseEntity.status(204).build();
    }
}
