package client

import (
	"context"
	"fmt"
	"github.com/metaprov/mdgoclient/gen/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"time"
)

// SDK is an instance of the Agones SDK
type PredictorClient struct {
	client api.PredictionServerClient
	ctx    context.Context
	host   string
	port   int32
}

func NewPredictorClient(host string, port int32) (*PredictorClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	s := &PredictorClient{
		ctx:  context.Background(),
		host: host,
		port: port,
	}
	// block for at least 30 seconds
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return s, errors.Wrapf(err, "could not connect to %s", addr)
	}
	s.client = api.NewPredictionServerClient(conn)
	return s, errors.Wrap(err, "could not set up health check")
}

func (r *PredictorClient) Ready() error {
	_, err := r.client.Ready(r.ctx, &api.ReadyRequest{})
	if err != nil {
		return errors.Wrap(err, "failed not send Ready message")
	}
	return nil

}

func (r *PredictorClient) Alive() error {
	_, err := r.client.Ready(r.ctx, &api.ReadyRequest{})
	if err != nil {
		return errors.Wrap(err, "failed not send Ready message")
	}
	return nil
}

func (r *PredictorClient) GetProduct() (string, error) {
	product, err := r.client.GetProduct(r.ctx, &api.GetProductRequest{})
	if err != nil {
		return "", errors.Wrap(err, "failed not send Ready message")
	}
	return product.Content, nil
}

func (r *PredictorClient) GetSchema() (string, error) {
	schema, err := r.client.GetSchema(r.ctx, &api.GetSchemaRequest{})
	if err != nil {
		return "", errors.Wrap(err, "failed not send Get Schema message")
	}
	return schema.Content, nil
}

func (r *PredictorClient) GetDataset() (string, error) {
	dataset, err := r.client.GetDataset(r.ctx, &api.GetDatasetRequest{})
	if err != nil {
		return "", errors.Wrap(err, "failed not send Get Dataset message")
	}
	return dataset.Content, nil
}

func (r *PredictorClient) GetModel() (string, error) {
	model, err := r.client.GetModel(r.ctx, &api.GetModelRequest{})
	if err != nil {
		return "", errors.Wrap(err, "failed not send Get Model message")
	}
	return model.Content, nil
}

func (r *PredictorClient) GetStats() (string, error) {
	stat, err := r.client.GetStat(r.ctx, &api.GetStatRequest{})
	if err != nil {
		return "", errors.Wrap(err, "failed not send Get Model message")
	}
	return stat.Content, nil
}

func (r *PredictorClient) Predict(colsJson string, dataJson string, full bool) (string, error) {
	req := &api.PredictionRequest{
		Columns:  colsJson,
		Features: dataJson,
		Full:     false,
	}

	result, err := r.client.Predict(r.ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "failed prediction")
	}
	return result.Labels, nil
}
