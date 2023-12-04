package webserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	cnt "github.com/arslab/lwnsimulator/controllers"
	"github.com/arslab/lwnsimulator/models"
	gw "github.com/arslab/lwnsimulator/simulator/components/gateway"
	"github.com/arslab/lwnsimulator/socket"
	_ "github.com/arslab/lwnsimulator/webserver/statik"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rakyll/statik/fs"
)

// WebServer type
type WebServer struct {
	Address      string
	Port         int
	Router       *gin.Engine
	ServerSocket *socketio.Server
}

var (
	simulatorController cnt.SimulatorController
	configuration       *models.ServerConfig
)

func NewWebServer(config *models.ServerConfig, controller cnt.SimulatorController) *WebServer {

	serverSocket := newServerSocket()

	configuration = config
	simulatorController = controller

	go func() {

		err := serverSocket.Serve()

		if err != nil {
			log.Fatal(err)
		}

	}()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true
	configCors.AllowHeaders = []string{"Origin", "Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers", "Content-type"}
	configCors.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS"}
	configCors.AllowCredentials = true
	router.Use(cors.New(configCors))

	router.Use(gin.Recovery())

	ws := WebServer{
		Address:      configuration.Address,
		Port:         configuration.Port,
		Router:       router,
		ServerSocket: serverSocket,
	}

	staticFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	staticGroup := router.Group("/dashboard")
	staticGroup.StaticFS("/", staticFS)
	//router.Use(static.Serve("/", staticFS))

	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/start", startSimulator)
		apiRoutes.GET("/stop", stopSimulator)
		apiRoutes.GET("/status", simulatorStatus)
		apiRoutes.GET("/bridge", getRemoteAddress)
		apiRoutes.GET("/gateways", getGateways)
		apiRoutes.POST("/del-gateway", deleteGateway)
		apiRoutes.POST("/add-gateway", addGateway)
		apiRoutes.POST("/up-gateway", updateGateway)
		apiRoutes.POST("/bridge/save", saveInfoBridge)
	}

	router.GET("/socket.io/*any", gin.WrapH(serverSocket))
	router.POST("/socket.io/*any", gin.WrapH(serverSocket))

	router.GET("/", func(context *gin.Context) { context.Redirect(http.StatusMovedPermanently, "/dashboard") })

	return &ws
}

func startSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, simulatorController.Run())
}

func stopSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, simulatorController.Stop())
}

func simulatorStatus(c *gin.Context) {
	c.JSON(http.StatusOK, simulatorController.Status())
}

func saveInfoBridge(c *gin.Context) {

	var ns models.AddressIP
	c.BindJSON(&ns)

	c.JSON(http.StatusOK, gin.H{"status": simulatorController.SaveBridgeAddress(ns)})
}

func getRemoteAddress(c *gin.Context) {
	c.JSON(http.StatusOK, simulatorController.GetBridgeAddress())
}

func getGateways(c *gin.Context) {

	gws := simulatorController.GetGateways()
	c.JSON(http.StatusOK, gws)
}

func addGateway(c *gin.Context) {

	var g gw.Gateway
	c.BindJSON(&g)

	code, id, err := simulatorController.AddGateway(&g)
	errString := fmt.Sprintf("%v", err)

	c.JSON(http.StatusOK, gin.H{"status": errString, "code": code, "id": id})

}

func updateGateway(c *gin.Context) {

	var g gw.Gateway
	c.BindJSON(&g)

	code, err := simulatorController.UpdateGateway(&g)
	errString := fmt.Sprintf("%v", err)

	c.JSON(http.StatusOK, gin.H{"status": errString, "code": code})

}

func deleteGateway(c *gin.Context) {

	Identifier := struct {
		Id int `json:"id"`
	}{}

	c.BindJSON(&Identifier)

	c.JSON(http.StatusOK, gin.H{"status": simulatorController.DeleteGateway(Identifier.Id)})

}

func newServerSocket() *socketio.Server {

	serverSocket := socketio.NewServer(nil)

	serverSocket.OnConnect("/", func(s socketio.Conn) error {

		log.Println("[WS]: Socket connected")

		s.SetContext("")
		simulatorController.AddWebSocket(&s)

		return nil

	})

	serverSocket.OnDisconnect("/", func(s socketio.Conn, reason string) {
		s.Close()
	})

	serverSocket.OnEvent("/", socket.EventToggleStateGateway, func(s socketio.Conn, Id int) {
		simulatorController.ToggleStateGateway(Id)
	})

	return serverSocket
}

func (ws *WebServer) Run() {

	log.Println("[WS]: Listen [", ws.Address+":"+strconv.Itoa(ws.Port), "]")

	err := ws.Router.Run(ws.Address + ":" + strconv.Itoa(ws.Port))
	if err != nil {
		log.Println("[WS] [ERROR]:", err.Error())
	}

}
