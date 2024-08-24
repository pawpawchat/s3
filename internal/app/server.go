package app

import (
	"context"
	"log"
	"log/slog"
	"net"
	"sync"

	s3cfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pawpawchat/s3/api/pb"
	"github.com/pawpawchat/s3/config"
	"github.com/pawpawchat/s3/internal/domain/model"
	"github.com/pawpawchat/s3/internal/infrastructure/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type mediaRepository interface {
	CreateMedia(context.Context, *model.Media) error
}

type s3Server struct {
	pb.UnimplementedS3Server
	remote  *s3.Client
	presign *s3.PresignClient
	mediaRp mediaRepository
}

func newS3Server(remote *s3.Client, mediaRp mediaRepository) *s3Server {
	return &s3Server{
		remote:  remote,
		presign: s3.NewPresignClient(remote),
		mediaRp: mediaRp,
	}
}

func Run(ctx context.Context, config *config.Config) {
	srv, lsn := settingUpServer(config)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Debug("grpc server setting up and running on", "addr", config.Env().ServerAddr)
		if err := srv.Serve(lsn); err != nil {
			slog.Error("grpc server", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		srv.GracefulStop()
		slog.Debug("grpc server has been gracefully shut down")
	}()

	wg.Wait()
}

func settingUpServer(config *config.Config) (*grpc.Server, net.Listener) {
	lsn, err := net.Listen("tcp", config.Env().ServerAddr)
	if err != nil {
		log.Fatal(err)
	}

	pg := postgresConn(config.Env().DatabaseURL)
	srv := grpc.NewServer()
	s3srv := newS3Server(s3ClientConn(), repository.NewMediaRepository(pg))

	pb.RegisterS3Server(srv, s3srv)
	reflection.Register(srv)

	return srv, lsn
}

func s3ClientConn() *s3.Client {
	cfg, err := s3cfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return s3.NewFromConfig(cfg)
}

func postgresConn(url string) *sqlx.DB {
	return sqlx.MustOpen("pgx", url)
}
