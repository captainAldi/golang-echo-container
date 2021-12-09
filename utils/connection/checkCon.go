package connection

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/go-redis/redis/v8"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/streadway/amqp"
)

func CheckMySQL(dbUser, dbPass, dbName, dbHost, dbPort string) string {

	statusDB := "Connected"

	// Open connection to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", dbUser, dbPass, dbHost, dbPort, dbName)

	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		statusDB = fmt.Sprintf("Error: %s", err)
	}

	return statusDB
}

func CheckPostgreSQL(dbUser, dbPass, dbName, dbHost, dbPort string) string {
	statusDB := "Connected"

	// Open Connection to Database
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		statusDB = fmt.Sprintf("Error: %s", err)
	}

	return statusDB

}

func CheckMongo(dbUser, dbPass, dbName, dbHost, dbPort, dbAuthSource string) string {
	statusMongo := "Connected"

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	credential := options.Credential{
		Username:   dbUser,
		Password:   dbPass,
		AuthSource: dbAuthSource,
	}

	mongoUrl := fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	clientOpts := options.Client().ApplyURI(mongoUrl)

	if dbUser != "" || dbPass != "" {

		clientOpts = clientOpts.SetAuth(credential)
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		statusMongo = fmt.Sprintf("Error: %s", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		statusMongo = fmt.Sprintf("Error: %s", err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	return statusMongo
}

func CheckRedis(redisHost, redisPort, redisUser, redisPass string) string {
	statusRedis := "Connected"

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPass, // no password set
		DB:       0,         // use default DB
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		statusRedis = fmt.Sprintf("Error: %s", err)
	}

	return statusRedis

}

func CheckRabbitMQ(rabbitMqHost, rabbitMqPort, rabbitMqUser, rabbitMqPass string) string {

	statusRabbitMQ := "Connected"

	rabbitMQurl := fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitMqUser, rabbitMqPass, rabbitMqHost, rabbitMqPort)

	conn, err := amqp.Dial(rabbitMQurl)
	if err != nil {
		statusRabbitMQ = fmt.Sprintf("Error: %s", err)
	}

	if conn != nil {
		defer conn.Close()
	}

	return statusRabbitMQ
}

func CheckImageVersion(k8sNamespace, k8sDeploymentName string) string {
	statusImageVersion := "vx.x.x"

	// Create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		statusImageVersion = fmt.Sprintf("Error: %s", err)
	}

	// Create out cluster config
	// config, err := clientcmd.BuildConfigFromFlags("", k8sKubeConfig)
	// if err != nil {
	// 	statusImageVersion = fmt.Sprintf("Error: %s", err)
	// }

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		statusImageVersion = fmt.Sprintf("Error: %s", err)
	}

	// Cek Deployment
	deployment, _ := clientset.AppsV1().Deployments(k8sNamespace).Get(context.TODO(), k8sDeploymentName, v1.GetOptions{})

	// Get Image Version
	statusImageVersion = deployment.Spec.Template.Spec.Containers[0].Image

	return statusImageVersion

}
