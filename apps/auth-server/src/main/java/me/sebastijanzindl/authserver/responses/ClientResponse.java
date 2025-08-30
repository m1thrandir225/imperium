package me.sebastijanzindl.authserver.responses;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import me.sebastijanzindl.authserver.model.Client;

import java.util.List;
import java.util.UUID;

public record ClientResponse (
    @JsonProperty("id")
     UUID id,
    @JsonProperty("client_name")
     String clientName,
    @JsonProperty("ip_address")
     String ipAddress
)
{
    public static ClientResponse from(Client client) {
        return new ClientResponse(
                client.getId(),
                client.getName(),
                client.getIpAddress()
        );
    }

    public static List<ClientResponse> fromList(List<Client> clients) {
        return clients.stream()
                .map(ClientResponse::from)
                .toList();
    }
}
