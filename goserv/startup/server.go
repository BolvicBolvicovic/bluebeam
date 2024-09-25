package startup

import (
	"context"
	"time"
	
//	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/goserve/arch/mongo"
	"github.com/unusualcodeorg/goserve/arch/network"
	"github.com/unusualcodeorg/goserve/arch/redis"
	"github.com/unusualcodeorg/goserve/config"
)

type Shutdown = func()

func Server() {
	env := config.NewEnv(".env", true)
	router, _, shutdown := create(env)
	defer shutdown()
	router.Start(env.ServerHost, env.ServerPort)
}

func create(env *config.Env) (network.Router, Module, Shutdown) {
	ctx		:= context.Background()
	dbConf	:= mongo.DbConfig {
		User:		env.DBUser,
		Pwd :		env.DBUserPwd,
		Host:		env.DBHost,
		Port:		env.DBPort,
		Name:		env.DBName,
		MinPoolSize:env.DBMinPoolSize,
		MaxPoolSize:env.DBMaxPoolSize,
		Timeout:	time.Duration(env.DBQueryTimeout) * time.Second,
	}

	db := mongo.NewDatabase(ctx, dbConf)
	db.Connect()

	//if env.GoMode != gin.TestMode { EnsureDbIndexes(db) }

	redisConf := redis.Config {
		Host:	env.RedisHost,
		Port:	env.RedisPort,
		Pwd :	env.RedisPwd,
		DB	:	env.RedisDB,
	}

	store := redis.NewStore(ctx, &redisConf)
	store.Connect()

	module := NewModule(ctx, env, db, store)

	router := network.NewRouter(env.GoMode)
	router.RegisterValidationParsers(network.CustomTagNameFunc())
	router.LoadRootMiddlewares(module.RootMiddlewares())
	router.LoadControllers(module.Controllers())

	shutdown := func() {
		db.Disconnect()
		store.Disconnect()
	}

	return router, module, shutdown
}
