package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/Bccp-Team/bccp-server/mysql"
	pb "github.com/Bccp-Team/bccp-server/proto/api"
	"github.com/Bccp-Team/bccp-server/runners"
	"github.com/Bccp-Team/bccp-server/scheduler"
)

func (*server) RunList(ctx context.Context, in *pb.Criteria) (*pb.Runs, error) {
	runs, err := mysql.Db.ListRuns(in.Filters, in.Limit, in.Offset)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	result := &pb.Runs{runs}

	return result, nil
}

func (*server) RunStat(ctx context.Context, in *pb.Criteria) (*pb.RunStats, error) {
	stats, err := mysql.Db.StatRun(in.Filters)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	return stats, nil
}

func (*server) RunGet(ctx context.Context, in *pb.Run) (*pb.Run, error) {
	run, err := mysql.Db.GetRun(in.Id)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	return run, nil
}

func (*server) RunStart(ctx context.Context, in *pb.Run) (*pb.Run, error) {
	id, err := mysql.Db.AddRun(in.RepoId, in.Batch, in.Priority)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	run, err := mysql.Db.GetRun(id)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	scheduler.DefaultScheduler.AddRun(run)

	return run, nil
}

func (*server) RunCancel(ctx context.Context, in *pb.Run) (*pb.Run, error) {
	run, err := mysql.Db.GetRun(in.Id)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	if run.RunnerId != 0 {
		runners.KillRun(run.RunnerId, run.Id)
	}

	err = mysql.Db.UpdateRunStatus(run.Id, "canceled")

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	run, err = mysql.Db.GetRun(in.Id)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	return run, nil
}
