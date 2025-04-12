package repository

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"news_dashboard/internal/config"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchRepository struct {
	esClient *elasticsearch.Client
}

func NewElasticsearchRepository(esClient *elasticsearch.Client) *ElasticsearchRepository {
	return &ElasticsearchRepository{
		esClient: esClient,
	}
}

func InitIndex(indexName string, cfg config.Config) error {
	// url := "http://elasticsearch:9200/" + indexName
	// url := "http://localhost:9200/" + indexName
	url := cfg.ElasticSearch.URL + "/" + indexName

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(nil))
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	indexJSON, err := ioutil.ReadFile("myindex.json")
	if err != nil {
		fmt.Println("Error reading index file:", err)
		return err
	}

	// 建立請求
	reqPut, err := http.NewRequest("PUT", url, bytes.NewBuffer(indexJSON))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// 設定請求header
	reqPut.Header.Set("Content-Type", "application/json")

	// 發送請求
	clientPut := http.Client{}
	res, err := clientPut.Do(reqPut)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer res.Body.Close()

	// 讀取回應
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	fmt.Println(string(body))
	return nil
}

// func (r *ElasticsearchRepository) SavePost(post Post) error {
// 	indexName := fmt.Sprintf("%s-%s", post.Type(), post.CreatedAt().Format("2006.01.02"))

// 	body, err := json.Marshal(post)
// 	if err != nil {
// 		return err
// 	}

// 	req := esapi.IndexRequest{
// 		Index:      indexName,
// 		DocumentID: fmt.Sprintf("%d", post.ID()),
// 		Body:       bytes.NewReader(body),
// 	}

// 	res, err := req.Do(context.Background(), r.esClient)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		return fmt.Errorf("failed to save post with id=%d, status code=%d", post.ID(), res.StatusCode())
// 	}

// 	log.Printf("post with id=%d has been saved to index=%s\n", post.ID(), indexName)

// 	return nil
// }

// func (r *ElasticsearchRepository) GetPostsByKeyword(keyword string, postType models.PostType) ([]models.Post, error) {
// 	indexName := fmt.Sprintf("%s-*", postType)

// 	query := map[string]interface{}{
// 		"query": map[string]interface{}{
// 			"multi_match": map[string]interface{}{
// 				"query":  keyword,
// 				"fields": []string{"title", "content"},
// 			},
// 		},
// 	}

// 	var buf bytes.Buffer
// 	if err := json.NewEncoder(&buf).Encode(query); err != nil {
// 		return nil, err
// 	}

// 	req := esapi.SearchRequest{
// 		Index: []string{indexName},
// 		Body:  bytes.NewReader(buf.Bytes()),
// 	}

// 	res, err := req.Do(context.Background(), r.esClient)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		return nil, fmt.Errorf("failed to search posts with keyword=%s, status code=%d", keyword, res.StatusCode())
// 	}

// 	var searchRes struct {
// 		Hits struct {
// 			Hits []struct {
// 				Source json.RawMessage `json:"_source"`
// 			} `json:"hits"`
// 		} `json:"hits"`
// 	}

// 	if err := json.NewDecoder(res.Body).Decode(&searchRes); err != nil {
// 		return nil, err
// 	}

// 	posts := []models.Post{}
// 	for _, hit := range searchRes.Hits.Hits {
// 		post, err := models.NewPostFromJSON(postType, hit.Source)
// 		if err != nil {
// 			log.Printf("failed to parse post: %v\n", err)
// 			continue
// 		}

// 		posts = append(posts, post)
// 	}

// 	log.Printf("found %d posts with keyword=%s in type=%s\n", len(posts), keyword, postType)

// 	return posts, nil
// }
