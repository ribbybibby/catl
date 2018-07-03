package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/prometheus/common/log"

	elastic "github.com/olivere/elastic"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	var (
		app          = kingpin.New("catl", "Writes a logstash elasticsearch query to stdout as if it were a logfile")
		index        = app.Flag("index", "Index pattern.").Short('i').Default("logstash-*").String()
		url          = app.Flag("url", "Logstash server URL.").Short('u').Default("http://localhost:9200").String()
		messageField = app.Flag("message-field", "Field to be returned").Short('m').Default("message").String()
		sortField    = app.Flag("sort-field", "Field to sort the results by").Short('s').Default("@timestamp").String()
		query        = app.Arg("query", "Elasticseach query string.").Required().String()
	)

	app.DefaultEnvars()
	app.HelpFlag.Short('h')
	app.Parse(os.Args[1:])

	client, err := elastic.NewClient(
		elastic.SetURL(*url),
	)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	queryString := elastic.NewQueryStringQuery(*query)
	scroll := client.Scroll().Index(*index).Query(queryString).Sort(*sortField, true).Size(10).Pretty(true)

	for {
		results, err := scroll.Do(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		for _, item := range results.Hits.Hits {
			var m interface{}
			err := json.Unmarshal(*item.Source, &m)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(m.(map[string]interface{})[*messageField].(string))
		}
	}
}
