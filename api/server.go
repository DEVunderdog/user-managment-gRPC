package api

import (
	"crypto/rsa"
	"fmt"

	database "github.com/DEVunderdog/user-management-gRPC/database/sqlc"
	"github.com/DEVunderdog/user-management-gRPC/pb"
	"github.com/DEVunderdog/user-management-gRPC/token"
	"github.com/DEVunderdog/user-management-gRPC/utils"
)

type Server struct {
	pb.UnimplementedUserManagementServer
	config utils.Config
	store database.Store
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store database.Store, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (*Server, error) {
	
	tokenMaker, err := token.NewJWTMaker(publicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}