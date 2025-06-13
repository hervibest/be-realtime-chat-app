package repository

import (
	"be-realtime-chat-app/services/chat-command-elastic-svc/internal/entity"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
)

type MessageRepository interface {
	Insert(message *entity.Message) error
}

type messageRepository struct {
	esClient *elasticsearch.Client
	index    string
}

func NewMessageRepository(esClient *elasticsearch.Client) MessageRepository {
	return &messageRepository{
		esClient: esClient,
		index:    "messages",
	}
}

func (r *messageRepository) Insert(message *entity.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      r.index,
		DocumentID: message.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), r.esClient)
	if err != nil {
		return fmt.Errorf("error inserting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response: %s", res.String())
	}

	return nil
}

func (r *messageRepository) FindManyByRoomID(roomID string, limit int) ([]*entity.Message, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": map[string]interface{}{
					"term": map[string]interface{}{
						"room_id": roomID,
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{"created_at": "desc"},
		},
		"size": limit,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := r.esClient.Search(
		r.esClient.Search.WithContext(context.Background()),
		r.esClient.Search.WithIndex(r.index),
		r.esClient.Search.WithBody(&buf),
		r.esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source entity.Message `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	messages := make([]*entity.Message, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		messages[i] = &hit.Source
	}

	return messages, nil
}

func (r *messageRepository) SoftDelete(id string) error {
	updateScript := `{
		"script": {
			"source": "ctx._source.deleted_at = params.deleted_at",
			"lang": "painless",
			"params": {
				"deleted_at": "%s"
			}
		}
	}`
	updateScript = fmt.Sprintf(updateScript, time.Now().Format(time.RFC3339Nano))

	req := esapi.UpdateRequest{
		Index:      r.index,
		DocumentID: id,
		Body:       strings.NewReader(updateScript),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), r.esClient)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response: %s", res.String())
	}

	return nil
}
