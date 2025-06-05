package me.sebastijanzindl.authserver.model;

import jakarta.persistence.*;

import java.util.UUID;

@Entity
@Table(name = "hosts")
public class Host {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private UUID id;
    @Column(unique = true, nullable = false)
    private String ipAddress;
    @Column(unique = true, nullable = false)
    private Integer port;
}
