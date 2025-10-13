package me.sebastijanzindl.authserver.service;

import jakarta.persistence.EntityNotFoundException;
import me.sebastijanzindl.authserver.dto.CreateHostDTO;
import me.sebastijanzindl.authserver.dto.UpdateHostDTO;
import me.sebastijanzindl.authserver.model.Host;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.HOST_STATUS;
import me.sebastijanzindl.authserver.repository.HostRepository;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;

import java.net.ProxySelector;
import java.net.http.HttpClient;
import java.net.http.HttpResponse;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
public class HostService {
    private final HostRepository hostRepository;
    private final WebClient webClient;

    public HostService(HostRepository hostRepository, WebClient.Builder builder) {
        this.hostRepository = hostRepository;
        this.webClient = builder.build();
    }

    public Host create(CreateHostDTO input, User owner) {
        Host host = new Host();
        host.setName(input.getName());
        host.setIpAddress(input.getIpAddress());
        host.setPort(input.getPort());
        host.setStatus(HOST_STATUS.AVAILABLE);
        host.setOwner(owner);

        return hostRepository.save(host);
    }

    public Host getHost(UUID id){
        return this.findById(id);
    }

    public Host getHostByName(String name) {
        return this.hostRepository.findByName(name).orElseThrow(() ->  new EntityNotFoundException("Host with name " + name + " not found"));
    }

    public Host getOrCreateHost(CreateHostDTO dto, User owner) {
        Optional<Host> existing = this.hostRepository.findByNameAndIpAddressAndPort(
                dto.getName(),
                dto.getIpAddress(),
                dto.getPort()
        );

        return existing.orElseGet(() -> this.create(dto, owner));
    }

    public String getPrograms(UUID id) {
        Host host = this.findById(id);
        if (host.getStatus() != HOST_STATUS.AVAILABLE) {
            throw new IllegalStateException("Host is not available");
        }
        String url = String.format("http://%s:%d/api/session/programs", host.getIpAddress(), host.getPort());

        return webClient
                .get()
                .uri(url)
                .retrieve()
                .bodyToMono(String.class)
                .block();
    }

    public Host update(UUID id, UpdateHostDTO input) {
        Host host = this.findById(id);

        host.setIpAddress(input.getIpAddress());
        host.setPort(input.getPort());
        host.setStatus(input.getStatus());

        return hostRepository.save(host);
    }

    public Host updateStatus(UUID id, HOST_STATUS status) {
        Host host = this.findById(id);
        host.setStatus(status);
        return hostRepository.save(host);
    }

    public Host delete(UUID id) {
        Host host = this.findById(id);

        hostRepository.delete(host);
        return host;
    }

    private Host findById(UUID id) {
        return hostRepository.findById(id).orElseThrow(() ->  new EntityNotFoundException("Host with id " + id + " not found"));
    }
}
