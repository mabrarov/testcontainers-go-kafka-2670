# Repro repo for testcontainers/testcontainers-go issue #2670

Reproducing https://github.com/testcontainers/testcontainers-go/issues/2670 (multiple tries can be required):

1. Environment
    1. Go SDK 1.23+
    1. Windows 11 23H2
    1. WSL
    1. Docker Desktop v4.41.2
1. Remove Kafka image used for the test
    ```shell
    docker rmi -f confluentinc/confluent-local:7.5.0 && docker system prune --force 
    ```
1. Restart Docker Desktop or (better) restart Windows
1. Run the test
    ```shell
    go test -v -count 1 ./...
    ```

Output of test when issue happens

```text
=== RUN   TestKafkaContainerStart
    testcontainers_kafka_test.go:14: starting container using image: confluentinc/confluent-local:7.5.0
2025/05/14 01:36:49 github.com/testcontainers/testcontainers-go - Connected to docker:
  Server Version: 28.1.1
  API Version: 1.46
  Operating System: Docker Desktop
  Total Memory: 7944 MB
  Testcontainers for Go Version: v0.32.0
  Resolved Docker Host: npipe:////./pipe/docker_engine
  Resolved Docker Socket Path: //var/run/docker.sock
  Test SessionID: 6fefc8bf9a56308aaa005b5df7fa4c74d1a963a9bc9b81b425b51ef4b73f9343
  Test ProcessID: ae68bdc4-392a-48f7-8f76-0d320dd9e505
2025/05/14 01:36:49 🐳 Creating container for image testcontainers/ryuk:0.7.0
2025/05/14 01:36:49 ✅ Container created: 7d149b1ba68b
2025/05/14 01:36:49 🐳 Starting container: 7d149b1ba68b
2025/05/14 01:36:49 ✅ Container started: 7d149b1ba68b
2025/05/14 01:36:49 ⏳ Waiting for container id 7d149b1ba68b image: testcontainers/ryuk:0.7.0. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms}
2025/05/14 01:36:49 🔔 Container is ready: 7d149b1ba68b
2025/05/14 01:36:49 Failed to get image auth for https://index.docker.io/v1/. Setting empty credentials for the image: confluentinc/confluent-local:7.5.0. Error is:credentials not found in native keychain
2025/05/14 01:37:38 🐳 Creating container for image confluentinc/confluent-local:7.5.0
2025/05/14 01:37:38 ✅ Container created: 478db2838910
2025/05/14 01:37:38 🐳 Starting container: 478db2838910
2025/05/14 01:38:38 ✅ Container started: 478db2838910
2025/05/14 01:38:38 container logs (port not found
context deadline exceeded):

    testcontainers_kafka_test.go:17: container start failed: failed to start container: port not found
        context deadline exceeded
--- FAIL: TestKafkaContainerStart (109.70s)
FAIL
FAIL    github.com/mabrarov/testcontainers-go-kafka-2670        109.794s
FAIL
```

Containers of the failed test look like

```text
$ docker ps -a
CONTAINER ID   IMAGE                                COMMAND                  CREATED              STATUS              PORTS                                         NAMES
478db2838910   confluentinc/confluent-local:7.5.0   "sh -c 'while [ ! -f…"   19 seconds ago       Up 18 seconds       8082/tcp, 9092/tcp, 0.0.0.0:62372->9093/tcp   thirsty_diffie
7d149b1ba68b   testcontainers/ryuk:0.7.0            "/bin/ryuk"              About a minute ago   Up About a minute   0.0.0.0:62338->8080/tcp                       reaper_6fefc8bf9a56308aaa005b5df7fa4c74d1a963a9bc9b81b425b51ef4b73f9343
$ docker logs -f 478db2838910
$ echo $?
0
```

The same issue in debugger

![debugger screenshot](debugger.png)

Please note that this issue is a sort of race condition - it depends on time b/w container started and
github.com/testcontainers/testcontainers-go/modules/kafka checked for the mapped port.
It makes this issue happening intermittently, so multiple (5-20) tries can be required to reproduce issue.
