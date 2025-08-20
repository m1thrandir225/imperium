package me.sebastijanzindl.authserver.service;

import me.sebastijanzindl.authserver.dto.CreateHostDTO;
import me.sebastijanzindl.authserver.dto.UpdateHostDTO;
import me.sebastijanzindl.authserver.model.Host;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.model.enums.HOST_STATUS;
import me.sebastijanzindl.authserver.repository.HostRepository;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class HostService {
    private final HostRepository hostRepository;

    public HostService(HostRepository hostRepository) {
        this.hostRepository = hostRepository;
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

    public Host getHost(UUID id) throws Exception {
        return hostRepository.findById(id).orElseThrow();
    }

    public Host update(UUID id, UpdateHostDTO input) {
        Host host = hostRepository.findById(id).orElseThrow();

        host.setIpAddress(input.getIpAddress());
        host.setPort(input.getPort());
        host.setStatus(input.getStatus());

        return hostRepository.save(host);
    }

    public Host updateStatus(UUID id, HOST_STATUS status) {
        Host host = hostRepository.findById(id).orElseThrow();
        host.setStatus(status);
        return hostRepository.save(host);
    }

    public Host delete(UUID id) {
        Host host = hostRepository.findById(id).orElseThrow();
        hostRepository.delete(host);
        return host;
    }
}
