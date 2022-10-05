package swarm

import (
	"context"
	"log"
	"net"

	pb "github.com/MindHunter86/aniliSeeder/swarm/grpc"
	"google.golang.org/grpc"
)

type Master struct {
	pb.UnimplementedMasterServer
}

func NewMaster() *Master {
	return &Master{}
}

func (*Master) InitialPhase(ctx context.Context, in *pb.MasterRequest) (*pb.MasterReply, error) {
	log.Printf("Received: %v", in.GetAccessKey())
	return &pb.MasterReply{Version: "Hello " + in.GetAccessKey()}, nil
}

func (*Master) Bootstrap() error {
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMasterServer(s, &Master{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return err
}
