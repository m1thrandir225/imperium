package me.sebastijanzindl.authserver.responses;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;
import me.sebastijanzindl.authserver.model.Client;
import me.sebastijanzindl.authserver.model.Host;

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
}
