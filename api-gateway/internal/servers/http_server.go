package servers

import (
	"log"
	"sync"
	"github.com/Abelova-Grupa/Mercypher/api/internal/websocket"
	"github.com/Abelova-Grupa/Mercypher/api/internal/domain"
	"github.com/Abelova-Grupa/Mercypher/api/internal/middleware"
	"github.com/gin-gonic/gin"
)

// HttpServer serves incoming http requests (excluding websocket) such as login and
// register.
//
// Note to self:	It could be more optimal to remove register and unregister channels,
//					and to define envelope messages for that purpose. Something that
//					should be tested in the future.
type HttpServer struct {
	router 		*gin.Engine					// HTTP Servers internal gin router
	wg 			*sync.WaitGroup				// Wait group that holds for HTTP server routine
	gwIn		chan *domain.Envelope		// Channel for sending envelopes to gateway
	gwOut		chan *domain.Envelope		// Channel for receiving envelopes from gateway
	register	chan *websocket.Websocket	// Channel for registering new user in gateway
	unregister	chan *websocket.Websocket	// Channel for unregistering user from gateway
}

func (s *HttpServer) handleLogin(ctx *gin.Context) {
	
}

func (s *HttpServer) handleRegister(ctx *gin.Context) {

}

func (s *HttpServer) handleLogout(ctx *gin.Context) {

}

func (s *HttpServer) handleWebSocket(ctx *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	ws := websocket.NewWebsocket(conn, domain.User{
		UserId: "example",
		Username: "testUser",
		Email: "test@user.rs",
	}, s.unregister)

	//TODO: Register this ws in gateway.
	s.register <- ws;


	// Handle this client in a new goroutine
	go ws.HandleClient()
}

func (s *HttpServer) setupRoutes() {
	
	// HTTP POST request routes
	//
	// Body of a login request should contain an username/email with password.
	// Body of a register request should contain a full user.
	s.router.POST("/login", s.handleLogin)
	s.router.POST("/register", s.handleRegister)

	// HTTP GET requset routes.
	//
	// Websocket route (/ws) must contain a valid token issued by login request.
	s.router.GET("/logout", s.handleLogout)
	s.router.GET("/ws", middleware.AuthMiddleware(), s.handleWebSocket)
} 

func NewHttpServer(wg *sync.WaitGroup, gwIn chan *domain.Envelope, gwOut chan *domain.Envelope, reg chan *websocket.Websocket, unreg chan *websocket.Websocket) *HttpServer {

	// Change to gin.DebugMode for development
	gin.SetMode(gin.ReleaseMode)

	server := &HttpServer{}
	router := gin.Default()

	// Server parameters
	server.wg = wg
	
	server.router = router
	server.setupRoutes()

	server.gwIn = gwIn
	server.gwOut = gwOut

	server.register = reg
	server.unregister = unreg

	return server
}

func (server *HttpServer) Start(address string) {
	server.wg.Add(1)
		
	// HTTP Server must run in its own routine for it has to work concurrently with
	// a gRPC server and main gateway router.
	go func() {
		defer server.wg.Done()
		
		log.Println("HTTP server thread started on: ", address)	
		if err := server.router.Run(address); err != nil {
			log.Fatal("Something went wrong while starting http server.")
		}
	}()	
}
