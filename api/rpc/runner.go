package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/Bccp-Team/bccp-server/mysql"
	pb "github.com/Bccp-Team/bccp-server/proto/api"
	"github.com/Bccp-Team/bccp-server/runners"
)

func (*server) RunnerList(ctx context.Context, in *pb.Criteria) (*pb.Runners, error) {
	runners := mysql.Db.ListRunners(in.Filters, in.Limit, in.Offset)

	result := &pb.Runners{runners}

	return result, nil
}

func (*server) RunnerGet(ctx context.Context, in *pb.Runner) (*pb.Runner, error) {
	runner, err := mysql.Db.GetRunner(in.Id)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	return runner, nil
}

func (*server) RunnerKill(ctx context.Context, in *pb.Runner) (*pb.Runner, error) {
	runners.KillRunner(in.Id)

	run, err := mysql.Db.GetRunner(in.Id)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	return run, nil
}

func (*server) RunnerStat(ctx context.Context, in *pb.Criteria) (*pb.RunnerStats, error) {
	stats, err := mysql.Db.StatRunners()

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	return stats, nil
}
