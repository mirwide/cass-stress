## Usage

    go build && ./cass-stress -h

    --seed value, -s value        specify cassandra node (default: "localhost:9042")
    --requests value, -r value    query count (default: 1000)
    --mode value, -m value        read/write
    --parallel value, -p value    goroutine count (default: 8)
    --connection value, -c value  connection per host (default: 8)
    --timeout value, -t value     connection timeout (default: 5s)
    --cql value                   cql version (default: 4)
    --replica-factor value        replica factor for keyspace (default: 1)
    --help, -h                    show help (default: false)