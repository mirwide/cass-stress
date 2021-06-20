package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"

	"github.com/gocql/gocql"
	"github.com/urfave/cli/v2"
)

func stress(ctx *cli.Context) error {
	node := ctx.String("seed")
	requests := ctx.Int("requests")
	mode := ctx.String("mode")
	goroutine := ctx.Int("parallel")
	replica := ctx.Int("replica-factor")

	cluster := gocql.NewCluster(node)
	cluster.ProtoVersion = ctx.Int("cql")
	cluster.NumConns = ctx.Int("connection")
	cluster.Timeout = ctx.Duration("timeout")
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// create keyspace
	createQuery := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS stress WITH replication = {'class': 'SimpleStrategy', 'replication_factor': %d}", replica)
	if err := session.Query(createQuery).Exec(); err != nil {
		log.Fatal(err)
	}

	// create table
	if err := session.Query("CREATE TABLE IF NOT EXISTS stress.stress_data ( key text , C0 blob, c1 blob, PRIMARY KEY(key))").Exec(); err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	g := 0
	start := 0
	for g < goroutine {
		g += 1
		wg.Add(1)
		end := start + requests/goroutine
		log.Printf("run thread %d, start token: %d, end token: %d", g, start, end)
		go func(i int) {
			defer wg.Done()
			for i < end {
				i += 1
				if mode == "write" {
					write(session, i)
				} else if mode == "read" {
					read(session, i)
				} else {
					log.Fatal("set mode to write or read")
				}
			}
		}(start)
		start = start + requests/goroutine
	}
	wg.Wait()
	return nil
}
func read(session *gocql.Session, i int) {
	key := fmt.Sprintf("stress_test_key_%d", i)
	var readKey string
	if err := session.Query("SELECT key, C0, C1 FROM stress.stress_data WHERE key = ?",
		key).Consistency(gocql.One).Scan(&readKey, nil, nil); err != nil {
		log.Print(err)
	}
	if readKey != key {
		log.Printf("problem read %s", key)
	}
}

func write(session *gocql.Session, i int) {
	key := fmt.Sprintf("stress_test_key_%d", i)
	n := 80
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)

	if err := session.Query("INSERT INTO stress.stress_data (key, C0, C1) VALUES (?, ?, ?)",
		key, s, s).Consistency(gocql.One).Exec(); err != nil {
		log.Print(err)
	}
}
