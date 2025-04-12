package main

import (
	"context"
	"fmt"
	"log"

	"news_dashboard/internal/config"
	"news_dashboard/internal/repository"
	"news_dashboard/parser"

	"github.com/olivere/elastic/v7"
)

func main() {

	// load config
	configFilePath := "config.yaml"
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	// create Document channel
	documentChan := make(chan *parser.Document)

	// Create parser factory
	factory := parser.NewParserFactory()

	// Create parsers for each source
	for _, parserItem := range cfg.Parsers {
		parser, err := factory.GetParser(parserItem, documentChan)
		if err != nil {
			log.Fatalf("Failed to create parser for %s: %s", parserItem.SrcType, err.Error())
		}
		// Start a goroutine to run parser
		go parser.Parse()
	}

	// for k := 1; k <= 100; k++ {
	// 	fmt.Println(k, <-documentChan)
	// }

	// Index Documents to Elasticsearch

	// create Elasticsearch client
	client, err := elastic.NewClient(
		elastic.SetURL(cfg.ElasticSearch.URL),
		elastic.SetSniff(false),
		// elastic.SetBasicAuth(viper.GetString("elastic_search.username"), viper.GetString("elastic_search.password")),
	)
	INDEX_NAME := "myindex"
	if err != nil {
		log.Fatalf("failed to create Elasticsearch client: %s", err)
	}

	err = repository.InitIndex(INDEX_NAME, cfg)
	if err != nil {
		log.Fatalf("failed to init Elasticsearch index: %s", err)
	}

	for doc := range documentChan {
		// indexReq := elastic.NewIndexRequest().Index(doc.Type).BodyJson(doc)
		// _, err := client.Index().Index(INDEX_NAME).Id(doc.ID).Doc(doc).Do(context.Background())
		_, err := client.Index().Index(INDEX_NAME).Id(doc.ID).BodyJson(doc).Do(context.Background())
		if err != nil {
			log.Fatalf("failed to index document with ID %s: %s", doc.ID, err.Error())
		} else {
			fmt.Printf("Document with ID %s is indexed successfully\n", doc.ID)
		}
	}
}
