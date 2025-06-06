package clients

import (
	"errors"

	sessionpb "github.com/Abelova-Grupa/Mercypher/session-service/external/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SessionClient struct {
	conn		*grpc.ClientConn
	client		sessionpb.SessionServiceClient
}

// NewSessionClient cretes a new client to a session service on the given address.
//
// Note:	The situation is the same as in NewMessageClient code. Even if the
//			connection fails or refuses it wont be registered. Only when sending
//			messages to an unexisting address will the error be thrown.
func NewSessionClient(address string) (*SessionClient, error){
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	if conn == nil {
		return nil, errors.New("Connection refused: nil")
	}

	client := sessionpb.NewSessionServiceClient(conn)

	return &SessionClient{
		conn: 	conn,
		client: client,
	}, nil
}

func (c *SessionClient) Close() error {
	return c.conn.Close()
}
