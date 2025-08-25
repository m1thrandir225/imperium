package me.sebastijanzindl.authserver.service;

import jakarta.persistence.EntityNotFoundException;
import me.sebastijanzindl.authserver.dto.UpdateUserDTO;
import me.sebastijanzindl.authserver.model.User;
import me.sebastijanzindl.authserver.repository.UserRepository;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
public class UserService {
    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;

    public UserService(UserRepository userRepository, PasswordEncoder passwordEncoder) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
    }

    public User update(UUID id, UpdateUserDTO input) {
        User me = this.findById(id);

        me.setName(input.getFirstName() + " " + input.getLastName());
        me.setEmail(input.getEmail());

        return this.userRepository.save(me);
    }

    public User updatePassword(UUID id, String oldPassword, String newPassword){
        User user = this.findById(id);

        if(!user.getPassword().equals(newPassword)) {
            throw new RuntimeException("Old password does not match");
        }

        String hashedPassword = passwordEncoder.encode(newPassword);

        user.setPassword(hashedPassword);

        return this.userRepository.save(user);
    }

    public User delete(UUID id) {
        User user = this.findById(id);

        this.userRepository.delete(user);
        return user;
    }

    private User findById(UUID id) {
        return this.userRepository.findById(id).orElseThrow(() ->  new EntityNotFoundException("User with id " + id + " not found"));
    }
}
