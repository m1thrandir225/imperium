package me.sebastijanzindl.authserver.model;

import jakarta.persistence.*;
import lombok.*;

import java.util.UUID;

@Getter
@Setter
@EqualsAndHashCode
@Entity
@Table(
        name = "clients",
        uniqueConstraints = {
                @UniqueConstraint(
                        name = "uk_client_name_ip_address",
                        columnNames = {"name", "ip_address"}
                )
        }
)
public class Client {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    @Column(columnDefinition = "uuid", updatable = false, nullable = false)
    private UUID id;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false, name = "ip_address")
    private String ipAddress;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "owner_id", nullable = false)
    private User owner;
}
