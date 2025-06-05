package me.sebastijanzindl.authserver.service;

import me.sebastijanzindl.authserver.repository.HostRepository;
import org.springframework.stereotype.Service;

@Service
public class HostService {
    private final HostRepository hostRepository;

    public HostService(HostRepository hostRepository) {
        this.hostRepository = hostRepository;
    }
}
