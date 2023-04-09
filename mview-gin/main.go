package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"log"
	"mview-gin/common"
	"net/http"
)

type Vote struct {
	Candidate string `json:"candidate"`
}

func main() {
	// 设置 Kafka 生产者
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	brokers := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating producer: %s", err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Error closing producer: %s", err)
		}
	}()

	// 设置 Kafka 消费者
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating consumer: %s", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing consumer: %s", err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("vote-events", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %s", err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalf("Error closing partition consumer: %s", err)
		}
	}()

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				var vote Vote
				if err := json.Unmarshal(msg.Value, &vote); err != nil {
					log.Printf("Error unmarshalling message: %s", err)
					continue
				}

				if _, ok := common.VoteCounts[vote.Candidate]; !ok {
					common.VoteCounts[vote.Candidate] = 0
				}
				common.VoteCounts[vote.Candidate] += 1
				log.Printf("Updated vote count: %+v", common.VoteCounts)

			case err := <-partitionConsumer.Errors():
				log.Printf("Error received from partition consumer: %s", err)
			}
		}
	}()

	// 设置 HTTP 服务
	r := gin.Default()

	r.POST("/vote", func(c *gin.Context) {
		var req Vote
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error111": err.Error()})
			return
		}

		jsonVote, err := json.Marshal(req)
		if err != nil {
			log.Printf("Error marshalling vote: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process vote"})
			return
		}

		msg := &sarama.ProducerMessage{
			Topic: "vote-events",
			Value: sarama.StringEncoder(jsonVote),
		}

		msg.Value = sarama.ByteEncoder(jsonVote)

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Error sending message: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send vote"})
			return
		}
		log.Printf("Sent message to partition %d at offset %d", partition, offset)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/results", func(c *gin.Context) {
		c.JSON(http.StatusOK, common.VoteCounts)
	})

	// 启动 HTTP 服务
	_ = r.Run(":3000")
}
