package me.sebastijanzindl.authserver.service;

import me.sebastijanzindl.authserver.dto.UpdateUserDTO;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.repository.UserRepository;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class UserService {
    private final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public User update(UUID id, UpdateUserDTO input) {
        User me = this.userRepository.findById(id).orElseThrow(() -> new RuntimeException("User with id " + id + " not found"));

        me.setName(input.getName());
        me.setEmail(input.getEmail());

        return this.userRepository.save(me);
    }

    public User updatePassword(UUID id, String oldPassword, String newPassword) {
        User user = this.userRepository.findById(id).orElseThrow(() -> new RuntimeException("User with id " + id + " not found"));

        user.setPassword(newPassword);
    }

    public User delete(UUID id) {}
}
