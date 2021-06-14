package utils

import (
	"context"
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

const ProjectID = "portfolio-tracker-210511"
const DatasetID = "binance"

var bigqueryClient *bigquery.Client

func getBigQueryClient() (client *bigquery.Client, err error) {

	if bigqueryClient == nil {
		bigqueryClient, err = bigquery.NewClient(context.Background(), ProjectID)
		if err != nil {
			return nil, err
		}
	}

	bigqueryClient.Location = "asia-east1"
	return bigqueryClient, nil
}

func QueryData(q string) ([]bigquery.Value, error) {

	client, err := getBigQueryClient()
	if err != nil {
		return nil, err
	}

	it, err := client.Query(q).Read(context.Background())
	if err != nil {
		return nil, err
	}

	values := []bigquery.Value{}
	for {
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func BqTableExists(tbName string) bool {
	
	client, err := getBigQueryClient()
	if err != nil {
		return false
	}

	table := client.Dataset(DatasetID).Table(tbName)
	if _, err := table.Metadata(context.Background()); err != nil {
		return false
	}

	return true
}

func CreateBqTable(tbName string, schema bigquery.Schema) error {

	client, err := getBigQueryClient()
	if err != nil {
		return err
	}

	if err := client.Dataset(DatasetID).Table(tbName).
		Create(context.Background(), &bigquery.TableMetadata{Schema: schema}); err != nil {
		return err
	}

	return nil
}

func InsertData(tbName string, data interface{}) error {

	client, err := getBigQueryClient()
	if err != nil {
		return err
	}

	uploader := client.Dataset(DatasetID).Table(tbName).Inserter()
	if err := uploader.Put(context.Background(), data); err != nil {
		return err
	}

	return nil
}