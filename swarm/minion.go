package swarm

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
)

type Minion struct{}

func NewMinion() *Minion {
	return &Minion{}
}

func (*Minion) Bootstrap() error {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMasterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.InitialPhase(ctx, &pb.MasterRequest{AccessKey: "fuckyounigga"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetVersion())
	return nil
}
