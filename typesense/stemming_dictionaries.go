package typesense

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type StemmingDictionariesInterface interface {
	Upsert(ctx context.Context, dictionaryId string, wordRootCombinations []api.StemmingDictionaryWord) ([]*api.StemmingDictionaryWord, error)
	UpsertJsonl(ctx context.Context, dictionaryId string, body io.Reader) (io.ReadCloser, error)
	Retrieve(ctx context.Context) (*api.ListStemmingDictionariesResponse, error)
}

type stemmingDictionaries struct {
	apiClient APIClientInterface
}

func (s *stemmingDictionaries) Upsert(ctx context.Context, dictionaryId string, wordRootCombinations []api.StemmingDictionaryWord) ([]*api.StemmingDictionaryWord, error) {
	var buf bytes.Buffer
	jsonEncoder := json.NewEncoder(&buf)
	for _, combo := range wordRootCombinations {
		if err := jsonEncoder.Encode(combo); err != nil {
			return nil, err
		}
	}

	response, err := s.UpsertJsonl(ctx, dictionaryId, &buf)
	if err != nil {
		return nil, err
	}

	var result []*api.StemmingDictionaryWord
	jsonDecoder := json.NewDecoder(response)
	for jsonDecoder.More() {
		var docResult *api.StemmingDictionaryWord
		if err := jsonDecoder.Decode(&docResult); err != nil {
			return result, err
		}
		result = append(result, docResult)
	}

	return result, nil
}

func (s *stemmingDictionaries) UpsertJsonl(ctx context.Context, dictionaryId string, body io.Reader) (io.ReadCloser, error) {
	params := &api.ImportStemmingDictionaryParams{
		Id: dictionaryId,
	}

	response, err := s.apiClient.ImportStemmingDictionaryWithBody(ctx,
		params, "application/octet-stream", body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		return nil, &HTTPError{Status: response.StatusCode, Body: body}
	}
	return response.Body, nil
}

func (s *stemmingDictionaries) Retrieve(ctx context.Context) (*api.ListStemmingDictionariesResponse, error) {
	response, err := s.apiClient.ListStemmingDictionariesWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		emptySlice := make([]string, 0)
		return &api.ListStemmingDictionariesResponse{
			JSON200: &struct {
				Dictionaries *[]string `json:"dictionaries,omitempty"`
			}{
				Dictionaries: &emptySlice,
			},
		}, nil
	}
	return response, nil
}
