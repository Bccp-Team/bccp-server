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

func (*server) NamespaceList(ctx context.Context, in *pb.Criteria) (*pb.Namespaces, error) {
	namespaces, err := mysql.Db.ListNamespaces()

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	return &pb.Namespaces{namespaces}, nil

}

func (*server) NamespaceGet(ctx context.Context, in *pb.Namespace) (*pb.Namespace, error) {
	namespace, err := mysql.Db.GetNamespace(in.Name)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	repos, err := mysql.Db.GetNamespaceRepos(&in.Name)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	namespace.Repos = repos

	return namespace, nil
}

func (*server) NamespaceCreate(ctx context.Context, in *pb.Namespace) (*pb.Namespace, error) {
	err := mysql.Db.AddNamespace(in.Name, in.IsCi)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	for _, repo := range in.Repos {
		_, err = mysql.Db.AddRepoToNamespace(in.Name, repo.Repo, repo.Ssh)
		if err != nil {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	repos, err := mysql.Db.GetNamespaceRepos(&in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	return &pb.Namespace{in.Name, repos, in.IsCi}, nil
}

func (*server) NamespaceAddRepo(ctx context.Context, in *pb.Namespace) (*pb.Namespace, error) {
	for _, repo := range in.Repos {
		_, err := mysql.Db.AddRepoToNamespace(in.Name, repo.Repo, repo.Ssh)
		if err != nil {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	repos, err := mysql.Db.GetNamespaceRepos(&in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	return &pb.Namespace{in.Name, repos, in.IsCi}, nil
}

func (*server) ReposDesactivate(ctx context.Context, in *pb.Repos) (*pb.Repos, error) {
	for _, repo := range in.Repos {
		err := mysql.Db.UpdateRepoActivation(repo.Id, false)
		if err != nil {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	return in, nil
}

func (*server) RepoPush(ctx context.Context, in *pb.Repo) (*pb.Runs, error) {
	repos, err := mysql.Db.GetCiReposFromName(in.Repo)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	var runs pb.Runs

	for _, repo := range repos {
		running, err := mysql.Db.ListRuns(map[string]string{
			"repo":   strconv.FormatInt(repo.Id, 10),
			"status": "waiting"},
			0, 0)

		if err != nil || len(running) > 0 {
			continue
		}

		batch, err := mysql.Db.GetLastBatchFromNamespace(repo.Namespace)

		if err != nil {
			//FIXME: log
			continue
		}

		runID, err := mysql.Db.AddRun(repo.Id, batch.Id, 5)

		if err != nil {
			//FIXME: log
			continue
		}

		run, err := mysql.Db.GetRun(runID)

		scheduler.DefaultScheduler.AddRun(run)

		runs.Runs = append(runs.Runs)
	}
	return &runs, nil
}

func (*server) NamespaceToggleCI(ctx context.Context, in *pb.Namespace) (*pb.Namespace, error) {
	err := mysql.Db.ToggleCI(*in)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	return in, nil
}
