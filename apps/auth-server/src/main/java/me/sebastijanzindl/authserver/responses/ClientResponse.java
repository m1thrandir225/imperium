package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import me.sebastijanzindl.authserver.model.Client;

import java.util.UUID;

@Data
@AllArgsConstructor
public class ClientResponse {
    @JsonProperty("id")
    public UUID id;

    @JsonProperty("client_name")
    public String clientName;

    @JsonProperty("ip_address")
    public String ipAddress;
    
    public ClientResponse(Client client) {
        this.id = client.getId();
        this.clientName = client.getName();
        this.ipAddress = client.getIpAddress();
    }
}
