package atom

import (
	"context"
	"encoding/json"
	"log"
	//	"strconv"
	"io"
	"strings"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func (atom *Atom) InitElasticsearch(addresses []string, username string, password string) error {
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	//  Get cluster info
	var r map[string]interface{}
	res, err := es.Info()
	if err != nil {
		log.Fatalf("ES, Error getting response: %s", err)
	}
	log.Println(res)
	// Check response status
	if res.IsError() {
		log.Fatalf("ES, Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("ES, Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("ES Client: %s", elasticsearch.Version)
	log.Printf("ES Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	atom.EsClient = es
	return nil
}
func (atom *Atom) SendFlow(index string, body io.Reader) error {
	// 1. Index document
	// Set up the request object directly.
	req := esapi.IndexRequest{
		Index:   index,
		Body:    body,
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), atom.EsClient)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	log.Println(res)
	if res.IsError() {
		log.Printf("[%s] Error indexing document: %s", res.Status(), req)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	log.Println(strings.Repeat("-", 37))

	// 2. Search for the indexed documents
	// Use the helper methods of the client.
	res, err = atom.EsClient.Search(
		atom.EsClient.Search.WithContext(context.Background()),
		atom.EsClient.Search.WithIndex("test"),
		atom.EsClient.Search.WithBody(strings.NewReader(`{"query" : { "match" : { "title" : "test" } }}`)),
		atom.EsClient.Search.WithTrackTotalHits(true),
		atom.EsClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))

	return nil
}

func EsTest() error {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://dev.regroupstudio.com/es",
		},
		Username: "elastic",
		Password: "regroupstudio",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	var r map[string]interface{}

	// 1. Get cluster info
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	log.Println(res)
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("ES Client: %s", elasticsearch.Version)
	log.Printf("ES Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	// 2. Index document
	// Set up the request object directly.
	title := "test atom"
	req := esapi.IndexRequest{
		Index: "test",
		// DocumentID: strconv.Itoa(i + 1),
		Body:    strings.NewReader(`{"title" : "` + title + `"}`),
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err = req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	log.Println(res)
	if res.IsError() {
		log.Printf("[%s] Error indexing document: %s", res.Status(), req)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	log.Println(strings.Repeat("-", 37))

	// 3. Search for the indexed documents
	// Use the helper methods of the client.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithBody(strings.NewReader(`{"query" : { "match" : { "title" : "test" } }}`)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))

	return nil
}
