package me.sebastijanzindl.authserver.model;

import jakarta.persistence.*;
import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.Setter;
import me.sebastijanzindl.authserver.model.enums.HOST_STATUS;

import java.util.UUID;

@Entity
@Getter
@Setter
@EqualsAndHashCode
@Table(
        name = "hosts",
        uniqueConstraints = {
                @UniqueConstraint(
                        name = "uk_host_name_ip_port_status",
                        columnNames = {"name", "ip_address", "port", "status"}
                )
        }
)
public class Host {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    @Column(columnDefinition = "uuid", updatable = false, nullable = false)
    private UUID id;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false)
    private String ipAddress;

    @Column(nullable = false)
    private Integer port;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, columnDefinition = "host_status_enum")
    private HOST_STATUS status = HOST_STATUS.AVAILABLE;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "owner_id", nullable = false)
    private User owner;
}
