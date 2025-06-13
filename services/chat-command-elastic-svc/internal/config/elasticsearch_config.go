package config

import (
	"be-realtime-chat-app/services/commoner/utils"

	"github.com/elastic/go-elasticsearch/v9"
)

func NewElasticsearch() (*elasticsearch.Client, error) {
	elasticAddres := utils.GetEnv("ELASTICSEARCH_ADDRESS")

	addresses := []string{
		elasticAddres,
	}

	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	res, err := client.Ping()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return client, nil
}

// func CreateMessageIndex(client *elasticsearch.Client) error {
// 	indexSettings := `{
// 		"settings": {
// 			"number_of_shards": 1,
// 			"number_of_replicas": 1
// 		},
// 		"mappings": {
// 			"properties": {
// 				"id": {"type": "keyword"},
// 				"uuid": {"type": "keyword"},
// 				"room_id": {"type": "keyword"},
// 				"user_id": {"type": "keyword"},
// 				"username": {"type": "keyword"},
// 				"content": {"type": "text"},
// 				"created_at": {"type": "date"},
// 				"deleted_at": {"type": "date"}
// 			}
// 		}
// 	}`

// 	res, err := client.Indices.Create(
// 		"messages",
// 		client.Indices.Create.WithBody(strings.NewReader(indexSettings)),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		return fmt.Errorf("error creating index: %s", res.String())
// 	}

// 	log.Println("Successfully created 'messages' index")
// 	return nil
// }
