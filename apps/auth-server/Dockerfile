FROM eclipse-temurin:21-jdk-alpine AS builder

WORKDIR /app

COPY gradlew .
COPY gradle/ gradle/
COPY build.gradle.kts .
COPY settings.gradle.kts . 

RUN chmod +x ./gradlew
RUN ./gradlew dependencies --no-daemon

COPY src/ src/

RUN ./gradlew build -x test --no-daemon

FROM eclipse-temurin:21-jre-alpine

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

RUN apk add --no-cache curl

WORKDIR /app

COPY --from=builder /app/build/libs/*.jar app.jar

RUN mkdir -p /app/logs && \
    chown -R appuser:appgroup /app

USER appuser

ENV JAVA_OPTS="-Xms512m -Xmx1024m -XX:+UseG1GC -XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0"

ENTRYPOINT [ "sh", "-c", "java ${JAVA_OPTS} -jar app.jar" ]