package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"golang.org/x/net/context"
	"io"
	"log"
	"time"
)

type elasticsearchClient struct {
	client *elasticsearch.Client
}

func NewElasticsearch() Search {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	cl, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalln("Elasticsearch error")
	}
	return &elasticsearchClient{client: cl}
}

func (e *elasticsearchClient) CreateIndex(indexName string) error {
	mappings := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type": "text",
				},
				"title_short": map[string]interface{}{
					"type": "text",
				},
				"description": map[string]interface{}{
					"type": "text",
				},
				"description_preview": map[string]interface{}{
					"type": "text",
				},
				"author": map[string]interface{}{
					"type": "keyword",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}

	body, err := json.Marshal(mappings)
	if err != nil {
		return err
	}

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("err: " + err.Error())
		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	log.Println("Index created successfully")
	return nil
}

func (e *elasticsearchClient) AddPost(post *models.Post) error {
	data, err := json.Marshal(*post)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "posts",
		DocumentID: fmt.Sprintf("%d", time.Now().UnixNano()),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("err: " + err.Error())
		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s", res.String())
	}

	log.Println("Indexed post successfully")
	return nil
}

func (e *elasticsearchClient) AddPostCategory(postCategory *models.PostCategory) error {
	data, err := json.Marshal(*postCategory)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "posts",
		DocumentID: fmt.Sprintf("%d", time.Now().UnixNano()),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("err: " + err.Error())
		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s", res.String())
	}

	log.Println("Indexed post successfully")
	return nil
}

func (e *elasticsearchClient) SearchPosts(query string) ([]models.Post, error) {
	var posts []models.Post

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title", "content", "author"},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return posts, err
	}

	req := esapi.SearchRequest{
		Index: []string{"posts"},
		Body:  bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return posts, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("err: " + err.Error())
		}
	}(res.Body)

	if res.IsError() {
		return posts, fmt.Errorf("error searching documents: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return posts, err
	}

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var post models.Post
		source := hit.(map[string]interface{})["_source"]
		sourceJSON, err := json.Marshal(source)
		if err != nil {
			return posts, err
		}
		if err := json.Unmarshal(sourceJSON, &post); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
