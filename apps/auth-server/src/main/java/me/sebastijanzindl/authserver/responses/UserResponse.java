package me.sebastijanzindl.authserver.responses;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;
import me.sebastijanzindl.authserver.model.Client;
import me.sebastijanzindl.authserver.model.Host;
import me.sebastijanzindl.authserver.model.User;

import java.util.Date;
import java.util.List;
import java.util.UUID;

@Getter
@AllArgsConstructor
public class UserResponse {
    private UUID id;
    private String email;
    private String name;
    private Date createdAt;
    private Date updatedAt;
    private List<Client> clients;
    private List<Host> hosts;

    public UserResponse(User user) {
        this.id = user.getId();
        this.email = user.getEmail();
        this.name = user.getName();
        this.createdAt = user.getCreatedAt();
        this.updatedAt = user.getUpdatedAt();
        this.clients = user.getClients();
        this.hosts = user.getHosts();
    }
}
