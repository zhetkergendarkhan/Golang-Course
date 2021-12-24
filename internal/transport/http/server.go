package http

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"net/http"
	"shop/internal/repository"
	"shop/internal/service"
	"shop/internal/transport/http/handler"
)

func initDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "final")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected!")
	return db, err
}

func StartServer(ch chan error) {
	db, err := initDB()
	if err != nil {
		ch <- err
		return
	}
	defer db.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use defgithub.com/ThreeDotsLabs/watermill-kafkaault DB
	})
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{"localhost:9092"},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		ch <- err
		return
	}
	saramaPublisherConfig := kafka.DefaultSaramaSyncPublisherConfig()
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:               []string{"localhost:9092"},
			Marshaler:             kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaPublisherConfig,
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		ch <- err
		return
	}
	pr := repository.NewProductRepositoryImpl(db)
	cr := repository.NewCategoryRepositoryImpl(db)

	ps := service.NewProductServiceImpl(pr, publisher)
	cs := service.NewCategoryServiceImpl(cr, rdb)

	prodHandler := handler.NewProductHandler(ps)
	catHandler := handler.NewCategoryHandler(cs)

	manager := handler.NewManager(prodHandler, catHandler)

	router := chi.NewRouter()
	messages, err := subscriber.Subscribe(context.Background(), "product-create")
	go manager.ProductCreate(messages)
	configureRouter(router, manager)

	ch <- http.ListenAndServe(":8000", router)

}
