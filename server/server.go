package server

import(
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/processor"
	"github.com/maxzerbini/ovo/server/model"
	"github.com/maxzerbini/ovo/command"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Server struct {
	keystorage storage.OvoStorage
	incmdproc *processor.InCommandQueue
	outcmdproc *processor.OutCommandQueue
	config *ServerConf	
	partitioner *processor.Partitioner
	innerServer *InnerServer
}

func NewServer(conf *ServerConf, ks storage.OvoStorage, in *processor.InCommandQueue, out *processor.OutCommandQueue) *Server {
	srv := &Server{keystorage:ks, incmdproc:in, outcmdproc:out, config:conf}
	srv.partitioner = processor.NewPartitioner(ks, &conf.ServerNode, out)
	srv.innerServer = NewInnerServer(conf, ks, in, out, srv.partitioner)
	return srv
}

func (srv *Server) Do() {
	go srv.innerServer.Do()
	// Creates a router without any middleware by default
    router := gin.New()
    // Global middleware
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
	router.GET("/ovo/keystorage/:key", srv.get )
	router.POST("/ovo/keystorage", srv.post )
	router.PUT("/ovo/keystorage", srv.post )
	router.DELETE("/ovo/keystorage", srv.delete )
	router.GET("/ovo/keystorage/:key/getandremove", srv.getAndRemove)
	router.POST("/ovo/keystorage/:key/updatevalueifequal", srv.updateValueIfEqual )
	router.POST("/ovo/keystorage/:key/updatekeyvalueifequal", srv.updateKeyAndValueIfEqual )
	router.POST("/ovo/keystorage/:key/updatekey", srv.updateKey )
	// Listen and server on 0.0.0.0:8080
	if srv.config.ServerNode.Node.Debug {
		gin.SetMode(gin.DebugMode)
	} else { gin.SetMode(gin.ReleaseMode) }
    router.Run(srv.config.ServerNode.Node.Host+":"+strconv.Itoa(srv.config.ServerNode.Node.Port))
}

func (srv *Server) get (c *gin.Context) {
	key := c.Param("key")
	if res,err := srv.keystorage.Get(key); err==nil {
		obj := model.NewOvoKVResponse(res)
		result := model.NewOvoResponse("done", "0", obj)
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, model.NewOvoResponse("error", "101", nil))
	}
}

func (srv *Server) post (c *gin.Context) {
	var kv model.OvoKVRequest
	if c.BindJSON(&kv) == nil {
		obj := model.NewMetaDataObj(&kv)
		srv.keystorage.Put(obj)
		srv.outcmdproc.Enqueu(&command.Command{OpCode:"put",Obj:obj.MetaDataUpdObj()})
		c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", nil))
	} else {
		c.JSON(http.StatusBadRequest, model.NewOvoResponse("error", "10", nil))
	}
}

func (srv *Server) delete (c *gin.Context) {
	key := c.Param("key")
	srv.keystorage.Delete(key);
	srv.outcmdproc.Enqueu(&command.Command{OpCode:"delete",Obj:&storage.MetaDataUpdObj{Key:key}})
	c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", nil))
}

func (srv *Server) getAndRemove (c *gin.Context) {
	key := c.Param("key")
	if res,err := srv.keystorage.GetAndRemove(key); err==nil {
		obj := model.NewOvoKVResponse(res)
		srv.outcmdproc.Enqueu(&command.Command{OpCode:"delete",Obj:&storage.MetaDataUpdObj{Key:key}})
		result := model.NewOvoResponse("done", "0", obj)
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusForbidden, model.NewOvoResponse("error", "102", nil))
	}
}

func (srv *Server) updateValueIfEqual (c *gin.Context) {
	key := c.Param("key")
	var kv model.OvoKVUpdateRequest
	if c.BindJSON(&kv) == nil {
		obj := model.NewMetaDataUpdObj(&kv)
		obj.Key = key
		err := srv.keystorage.UpdateValueIfEqual(obj)
		if err == nil {
			srv.outcmdproc.Enqueu(&command.Command{OpCode:"updatevalue",Obj:obj})
			c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", nil))
		} else {
			c.JSON(http.StatusForbidden, model.NewOvoResponse("error", "103", nil))
		}
	} else {
		c.JSON(http.StatusBadRequest, model.NewOvoResponse("error", "10", nil))
	}
}

func (srv *Server) updateKeyAndValueIfEqual (c *gin.Context) {
	key := c.Param("key")
	var kv model.OvoKVUpdateRequest
	if c.BindJSON(&kv) == nil {
		obj := model.NewMetaDataUpdObj(&kv)
		obj.Key = key
		err := srv.keystorage.UpdateKeyAndValueIfEqual(obj)
		if err == nil {
			srv.outcmdproc.Enqueu(&command.Command{OpCode:"updatekeyvalue",Obj:obj})
			c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", nil))
		} else {
			c.JSON(http.StatusForbidden, model.NewOvoResponse("error", "104", nil))
		}
	} else {
		c.JSON(http.StatusBadRequest, model.NewOvoResponse("error", "10", nil))
	}
}

func (srv *Server) updateKey(c *gin.Context) {
	key := c.Param("key")
	var kv model.OvoKVUpdateRequest
	if c.BindJSON(&kv) == nil {
		obj := model.NewMetaDataUpdObj(&kv)
		obj.Key = key
		err := srv.keystorage.UpdateKey(obj)
		if err == nil {
			srv.outcmdproc.Enqueu(&command.Command{OpCode:"updatekey",Obj:obj})
			c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", nil))
		} else {
			c.JSON(http.StatusForbidden, model.NewOvoResponse("error", "105", nil))
		}
	} else {
		c.JSON(http.StatusBadRequest, model.NewOvoResponse("error", "10", nil))
	}
}
	