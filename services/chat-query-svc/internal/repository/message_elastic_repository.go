package repository

import (
	"be-realtime-chat-app/services/chat-query-svc/internal/entity"
	"be-realtime-chat-app/services/chat-query-svc/internal/model"
	"bytes"
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/elastic/go-elasticsearch/v9"
)

type MessageElasticRepo interface {
	SearchMessages(params *model.SearchParams) (*[]*entity.Message, error)
}

type messageElasticRepoImpl struct {
	esClient *elasticsearch.Client
	index    string
}

func NewMessageElasticRepo(esClient *elasticsearch.Client) MessageElasticRepo {
	return &messageElasticRepoImpl{
		esClient: esClient,
		index:    "messages",
	}
}

func (r *messageElasticRepoImpl) SearchMessages(params *model.SearchParams) (*[]*entity.Message, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{},
		"sort": []map[string]interface{}{
			{"created_at": "desc"},
		},
		"size": params.Limit,
	}

	boolQuery := map[string]interface{}{}

	if params.Username != "" {
		boolQuery["must"] = append(boolQuery["must"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{"username": params.Username},
			})
	}

	if params.Content != "" {
		boolQuery["must"] = append(boolQuery["must"].([]map[string]interface{}),
			map[string]interface{}{
				"match": map[string]interface{}{"content": params.Content},
			})
	}

	if params.RoomID != "" {
		boolQuery["filter"] = append(boolQuery["filter"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{"room_id": params.RoomID},
			})
	}

	if len(boolQuery) > 0 {
		query["query"].(map[string]interface{})["bool"] = boolQuery
	}

	var buf bytes.Buffer
	if err := sonic.ConfigFastest.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query with sonic: %w", err)
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

	if err := sonic.ConfigFastest.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing the response body with sonic: %w", err)
	}

	messages := make([]*entity.Message, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		messages[i] = &hit.Source
	}

	return &messages, nil
}
