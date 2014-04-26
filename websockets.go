package main

import (
	"code.google.com/p/go.net/websocket"
	//"code.google.com/p/google-api-go-client/mirror/v1"
	"github.com/wearscript/wearscript-go/wearscript"
	"fmt"
	//"strings"
	"time"
	"sync"
)

var Locks = map[string]*sync.Mutex{} // [user]
var UserSockets = map[string][]chan *[]interface{}{} // [user]
var Managers = map[string]*wearscript.ConnectionManager{}

func CurTime() float64 {
	return float64(time.Now().UnixNano()) / 1000000000.
}

func WSHandler(ws *websocket.Conn) {
	defer ws.Close()
	fmt.Println("Connected with glass")
	fmt.Println(ws.Request().URL.Path)
	userId, err := userID(ws.Request())
    //    if err != nil || userId == "" {
    //            path := strings.Split(ws.Request().URL.Path, "/")
	//	fmt.Println(path)
    //            if len(path) != 3 {
    //                    fmt.Println("Bad path")
    //                    return
    //            }
    //            userId, err = getSecretUser("ws", secretHash(path[len(path)-1]))
    //            if err != nil {
    //                    fmt.Println(err)
    //                    return
    //            }
    //    }
	//userId = SanitizeUserId(userId)
	//, err = mirror.New(authTransport(userId).Client()) // svc
	//if err != nil {
	//	LogPrintf("ws: mirror")
	//	//WSSend(c, &[]interface{}{[]uint8("error"), []uint8("Unable to create mirror transport")})
	//	return
	//}
	// TODO(brandyn): Add lock here when adding users
	if Managers[userId] == nil {
		Managers[userId], err = wearscript.ConnectionManagerFactory("server", "demo")
		if err != nil {
			LogPrintf("ws: cm")
			return
		}
        Managers[userId].Subscribe("image", func (c string, dataBin []byte, data []interface{}) {
            timestamp := data[1].(float64)
            image := data[2].(string)
            fmt.Println(fmt.Sprintf("Image[%s] Time[%f] Bytes[%d]", c, timestamp,           len(image)))
        })
  
        Managers[userId].Subscribe("sensors", func (c string, dataBin []byte, data []interface{}) {
            names := data[1]
            samples := data[2]
            fmt.Println(fmt.Sprintf("Sensors[%s]\n Names[%v]\n Samples[%v]\n", c, names,    samples))
        })
		//Managers[userId].SubscribeTestHandler()
		//Managers[userId].Subscribe("gist", func(channel string, dataRaw []byte, data []interface{}) {
		//	GithubGistHandle(Managers[userId], userId, data)
		//})
		//Managers[userId].Subscribe("weariverse", func(channel string, dataRaw []byte, data []interface{}) {
		//	WeariverseGistHandle(Managers[userId], userId, data)
		//})
	}
	cm := Managers[userId]
	conn, err := cm.NewConnection(ws) // con
	if err != nil {
		fmt.Println("Failed to create conn")
		return
	}
	cm.HandlerLoop(conn)
}
