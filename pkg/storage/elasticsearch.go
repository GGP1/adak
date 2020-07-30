package storage

import (
	"context"
	"fmt"

	elastic "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
)

var hitsMapping = `
{
	"settings": {
    	"number_of_shards": 1,
    	"number_of_replicas": 1
	},
   "mappings": {
       "properties": {
         "id": {
               "type": "text"
         },
         "footprint": {
               "type": "text"      
         },
         "path": {
               "type": "text"
         },
         "url": {
               "type": "text"
         },
         "language": {
               "type": "text"
         },
         "user_agent": {
               "type": "text"
         },
         "referer": {
               "type": "text"
         },
         "date": {
                "type": "date"
         }
     }
   }
}`

// Elasticsearch initializes the elasticsearch client
func Elasticsearch(ctx context.Context) (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't initialize elasticsearch")
	}

	err = indexExists(ctx, client, "hits", hitsMapping)
	if err != nil {
		return nil, err
	}

	fmt.Println("Elastisearch initialized on port 9200")

	return client, err
}

// indexExists checks if the index is already created, if not, creates it
func indexExists(ctx context.Context, client *elastic.Client, index, mapping string) error {
	exist, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't find the indices")
	}

	if !exist {
		_, err := client.CreateIndex(index).Body(mapping).Do(ctx)
		if err != nil {
			return errors.Wrap(err, "couldn't create the index")
		}
	}

	return nil
}
