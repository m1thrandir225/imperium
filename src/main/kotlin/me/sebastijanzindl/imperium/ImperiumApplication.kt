package me.sebastijanzindl.imperium

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class ImperiumApplication

fun main(args: Array<String>) {
    runApplication<ImperiumApplication>(*args)
}
