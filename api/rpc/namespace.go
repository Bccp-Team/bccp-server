package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/Bccp-Team/bccp-server/mysql"
	pb "github.com/Bccp-Team/bccp-server/proto/api"
)

func (*server) NamespaceList(ctx context.Context, in *pb.Criteria) (*pb.Namespaces, error) {
	namespaces, err := mysql.Db.ListNamespaces()

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	return &pb.Namespaces{namespaces}, nil

}

func (*server) NamespaceGet(ctx context.Context, in *pb.Namespace) (*pb.Namespace, error) {
	repos, err := mysql.Db.GetNamespaceRepos(&in.Name)

	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	return &pb.Namespace{in.Name, repos}, nil
}

func (*server) NamespaceCreate(ctx context.Context, in *pb.Namespace) (*pb.Namespace, error) {
	err := mysql.Db.AddNamespace(in.Name)
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
	return &pb.Namespace{in.Name, repos}, nil
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
	return &pb.Namespace{in.Name, repos}, nil
}
