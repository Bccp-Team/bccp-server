package rpc

import (
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/Bccp-Team/bccp-server/mysql"
	pb "github.com/Bccp-Team/bccp-server/proto/api"
	"github.com/Bccp-Team/bccp-server/scheduler"
)

func (*server) BatchStart(ctx context.Context, in *pb.Batch) (*pb.Runs, error) {
	batchID, err := mysql.Db.AddBatch(in.Namespace,
		in.InitScript,
		in.UpdateTime,
		in.Timeout)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	repos, err := mysql.Db.GetNamespaceRepos(&in.Namespace)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	for _, repo := range repos {
		runID, err := mysql.Db.AddRun(int64(repo.ID), batchID)

		if err != nil {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}

		scheduler.DefaultScheduler.AddRun(runID)
	}

	runs, err := mysql.Db.ListRuns(map[string]string{"batch": strconv.FormatInt(batchID, 10)}, 0, 0)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	return &pb.Runs{runs}, nil
}
