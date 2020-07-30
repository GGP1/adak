// Package tracking is a privacy-focused user tracker inspired on
// github.com/emvi/pirsch
package tracking

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	elastic "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
)

// Hitter is the interface that wraps Hit methods.
type Hitter interface {
	Hit(ctx context.Context, r *http.Request) error
	DeleteHit(ctx context.Context, id string) error
	Get(ctx context.Context) ([]Hit, error)
}

// Searcher is the interface that serves searching methods.
type Searcher interface {
	Search(ctx context.Context, value string) ([]interface{}, error)
}

// Tracker provides methods to track requests and store them in a data store.
// In case of an error it will panic.
type Tracker struct {
	ESClient *elastic.Client
	salt     string
}

// NewTracker returns a new user tracker.
func NewTracker(esClient *elastic.Client, salt string) *Tracker {
	return &Tracker{
		ESClient: esClient,
		salt:     salt,
	}
}

// DeleteHit takes away the hit with the id specified from the database.
func (t *Tracker) DeleteHit(ctx context.Context, id string) error {
	bq := elastic.NewBoolQuery()
	bq.Must(elastic.NewTermQuery("id", id))
	_, err := elastic.NewDeleteByQueryService(t.ESClient).
		Index("hits").
		Query(bq).
		Do(ctx)

	if err != nil {
		return errors.Wrap(err, "couldn't delete the hit")
	}

	return nil
}

// Get lists all the hits stored in the database.
func (t *Tracker) Get(ctx context.Context) ([]Hit, error) {
	var hits []Hit

	searchSrc := elastic.NewSearchSource()
	searchSrc.Query(elastic.NewMatchAllQuery())

	searchSv := t.ESClient.Search().Index("hits").SearchSource(searchSrc)

	searchResult, err := searchSv.Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find the hit")
	}

	for _, r := range searchResult.Hits.Hits {
		var hit Hit
		err := json.Unmarshal(r.Source, &hit)
		if err != nil {
			return nil, errors.Wrap(err, "failed unmarshaling data")
		}

		hits = append(hits, hit)
	}

	return hits, nil
}

// Hit stores the given request.
// The request might be ignored if it meets certain conditions.
func (t *Tracker) Hit(ctx context.Context, r *http.Request) error {
	if !ignoreHit(r) {
		hit := HitRequest(r, t.salt)

		_, err := t.ESClient.Index().
			Index("hits").
			BodyJson(hit.String()).
			Do(ctx)

		if err != nil {
			return errors.Wrap(err, "couldn't save the hit")
		}
	}

	return nil
}

// Search looks for a value and returns a slice of the hits that contain that value.
func (t *Tracker) Search(ctx context.Context, value string) ([]Hit, error) {
	var hits []Hit

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewQueryStringQuery(value))

	searchService := t.ESClient.Search().Index("hits").SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "search failed")
	}

	for _, r := range searchResult.Hits.Hits {
		var hit Hit
		err := json.Unmarshal(r.Source, &hit)
		if err != nil {
			return nil, errors.Wrap(err, "failed unmarshaling data")
		}

		hits = append(hits, hit)
	}

	return hits, nil
}

// Check headers commonly used by bots.
// If the user is a bot return true, else return false.
func ignoreHit(r *http.Request) bool {
	if r.Header.Get("X-Moz") == "prefetch" ||
		r.Header.Get("X-Purpose") == "prefetch" ||
		r.Header.Get("X-Purpose") == "preview" ||
		r.Header.Get("Purpose") == "prefetch" ||
		r.Header.Get("Purpose") == "preview" {
		return true
	}

	userAgent := strings.ToLower(r.Header.Get("User-Agent"))

	for _, botUserAgent := range userAgentBotlist {
		if strings.Contains(userAgent, botUserAgent) {
			return true
		}

	}

	return false
}
