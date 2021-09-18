package atom

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Atom struct {
	GinEngine   *gin.Engine
	MongoClient *mongo.Client
	// RedisClient *redis.Client
	RedisClient *redis.ClusterClient
	EsClient    *elasticsearch.Client
	Log         *zap.SugaredLogger
	Viper       *viper.Viper
	MysqlDB     *gorm.DB
}

var GinEngine *gin.Engine
var MongoClient *mongo.Client
var RedisClient *redis.ClusterClient
var EsClient *elasticsearch.Client
var Log *zap.SugaredLogger
var Viper *viper.Viper
var MysqlDB *gorm.DB

func InitFrom(atom *Atom) {
	GinEngine = atom.GinEngine
	MongoClient = atom.MongoClient
	RedisClient = atom.RedisClient
	EsClient = atom.EsClient
	Log = atom.Log
	Viper = atom.Viper
	MysqlDB = atom.MysqlDB
}

type Router func(*Atom)
