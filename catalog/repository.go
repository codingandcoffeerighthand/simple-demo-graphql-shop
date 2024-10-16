package catalog

import (
	"context"
	"encoding/json"
	"errors"

	elastic "gopkg.in/olivere/elastic.v5"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var (
	ErrNotFound = errors.New("entity not found")
)

type Repository interface {
	Close()
	CreateProduct(ctx context.Context, product Product) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	GetListProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	GetListProductWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}
type productDocument struct {
	Name        string
	Description string
	Price       float64
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) CreateProduct(ctx context.Context, product Product) error {
	_, err := r.client.Index().
		Index("catalog").
		Type("product").
		Id(product.ID).
		BodyJson(productDocument{product.Name, product.Description, product.Price}).
		Do(ctx)
	return err
}

// Close implements Repository.
func (r *elasticRepository) Close() {
	panic("unimplemented")
}

// GetListProductWithIDs implements Repository.
func (r *elasticRepository) GetListProductWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().Index("catalog").Type("product").Id(id))
	}
	res, err := r.client.MultiGet().Add(items...).Do(ctx)
	if err != nil {
		return nil, err
	}
	var products []Product
	for _, doc := range res.Docs {
		p := productDocument{}
		if err := json.Unmarshal(*doc.Source, &p); err != nil {
			return nil, err
		}
		products = append(products, Product{
			ID:          doc.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}

// GetListProducts implements Repository.
func (r *elasticRepository) GetListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search().Index("catalog").Type("product").Query(elastic.NewMatchAllQuery()).From(int(skip)).Size(int(take)).Do(ctx)
	if err != nil {
		return nil, err
	}
	var products []Product
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(*hit.Source, &p); err != nil {
			return nil, err
		}
		products = append(products, Product{
			ID:          hit.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}

// GetProductById implements Repository.
func (r *elasticRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(ctx)

	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, ErrNotFound
	}
	p := productDocument{}
	if err := json.Unmarshal(*res.Source, &p); err != nil {
		return nil, err
	}
	return &Product{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price}, nil
}

// SearchProducts implements Repository.
func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(*hit.Source, &p); err != nil {
			return nil, err
		}
		products = append(products, Product{
			ID:          hit.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}
