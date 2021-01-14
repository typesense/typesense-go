# typesense-go 

[![Build Status](https://cloud.drone.io/api/badges/typesense/typesense-go/status.svg)](https://cloud.drone.io/typesense/typesense-go)
[![GoReportCard Status](https://goreportcard.com/badge/github.com/typesense/typesense-go)](https://goreportcard.com/report/github.com/typesense/typesense-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/typesense/typesense-go.svg)](https://pkg.go.dev/github.com/typesense/typesense-go)
<!-- ![GitHub release](https://img.shields.io/github/v/release/typesense/typesense-go) -->
<!-- [![Gitter](https://badges.gitter.im/typesense-go/community.svg)](https://gitter.im/typesense-go/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) -->


Go client for the Typesense API: https://github.com/typesense/typesense

## Installation

```
go get github.com/typesense/typesense-go
```

## Usage

Create new client:

```go
client := typesense.NewClient(
	    typesense.WithServer("http://localhost:8108"),
	    typesense.WithAPIKey("<API_KEY>"))
```

New client with advanced configuration options (see godoc):

```go
client := typesense.NewClient(
		typesense.WithServer("http://localhost:8108"),
		typesense.WithAPIKey("<API_KEY>"),
		typesense.WithConnectionTimeout(5*time.Second),
		typesense.WithCircuitBreakerMaxRequests(50),
		typesense.WithCircuitBreakerInterval(2*time.Minute),
		typesense.WithCircuitBreakerTimeout(1*time.Minute),
	)
```

You can also find some examples in [integration tests](https://github.com/typesense/typesense-go/tree/master/typesense/test).

### Create a collection

```go
	schema := &api.CollectionSchema{
		Name: "companies",
		Fields: []api.Field{
			{
				Name: "company_name",
				Type: "string",
			},
			{
				Name: "num_employees",
				Type: "int32",
			},
			{
				Name:  "country",
				Type:  "string",
				Facet: true,
			},
		},
		DefaultSortingField: "num_employees",
	}

	client.Collections().Create(schema)
```

### Index a document

```go
	document := struct {
		ID           string `json:"id"`
		CompanyName  string `json:"company_name"`
		NumEmployees int    `json:"num_employees"`
		Country      string `json:"country"`
	}{
		ID:           "123",
		CompanyName:  "Stark Industries",
		NumEmployees: 5215,
		Country:      "USA",
	}

	client.Collection("companies").Documents().Create(document)
```

### Upserting a document

```go
	document := struct {
		ID           string `json:"id"`
		CompanyName  string `json:"company_name"`
		NumEmployees int    `json:"num_employees"`
		Country      string `json:"country"`
	}{
		ID:           "123",
		CompanyName:  "Stark Industries",
		NumEmployees: 5215,
		Country:      "USA",
	}

	client.Collection("companies").Documents().Upsert(newDocument)
```

### Search a collection

```go
	searchParameters := &api.SearchCollectionParams{
		Q:        "stark",
		QueryBy:  []string{"company_name"},
		FilterBy: pointer.String("num_employees:>100"),
		SortBy:   &([]string{"num_employees:desc"}),
	}

	client.Collection("companies").Documents().Search(searchParameters)
```

### Retrieve a document

```go
client.Collection("companies").Document("123").Retrieve()
```

### Update a document

```go
	document := struct {
		CompanyName  string `json:"company_name"`
		NumEmployees int    `json:"num_employees"`
	}{
		CompanyName:  "Stark Industries",
		NumEmployees: 5500,
	}

	client.Collection("companies").Document("123").Update(document)
```

### Delete an individual document

```go
client.Collection("companies").Document("123").Delete()
```

### Delete a bunch of documents

```go
filter := &api.DeleteDocumentsParams{FilterBy: "num_employees:>100", BatchSize: 100}
client.Collection("companies").Documents().Delete(filter)
```

### Retrieve a collection

```go
client.Collection("companies").Retrieve()
```

### Export documents from a collection

```go
client.Collection("companies").Documents().Export()
```

### Import documents into a collection

The documents to be imported can be either an array of document objects or be formatted as a newline delimited JSON string (see [JSONL](https://jsonlines.org)).


Import an array of documents:

```go	
	documents := []interface{}{
		struct {
			ID           string `json:"id"`
			CompanyName  string `json:"companyName"`
			NumEmployees int    `json:"numEmployees"`
			Country      string `json:"country"`
		}{
			ID:           "123",
			CompanyName:  "Stark Industries",
			NumEmployees: 5215,
			Country:      "USA",
		},
	}
	params := &api.ImportDocumentsParams{
		Action:    "create",
		BatchSize: 40,
	}

	client.Collection("companies").Documents().Import(documents, params)
```

Import a JSONL file:

```go
	params := &api.ImportDocumentsParams{
		Action:    "create",
		BatchSize: 40,
	}
	importBody, err := os.Open("documents.jsonl")
	// defer close, error handling ...

	client.Collection("companies").Documents().ImportJsonl(importBody, params)
```

### List all collections

```go
client.Collections().Retrieve()
```

### Drop a collection

```go
client.Collection("companies").Delete()
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/typesense/typesense-go.

## License

`typesense-go` is distributed under the Apache 2 license.