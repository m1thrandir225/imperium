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
@Table(name = "hosts")
public class Host {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    @Column(columnDefinition = "uuid", updatable = false, nullable = false)
    private UUID id;
    @Column(unique = true, nullable = false)
    private String ipAddress;
    @Column(unique = true, nullable = false)
    private Integer port;

    @Column(unique = true, nullable = false)
    private HOST_STATUS status = HOST_STATUS.AVAILABLE;

    @ManyToOne(fetch = FetchType.LAZY, cascade = CascadeType.ALL)
    private User owner;
}
