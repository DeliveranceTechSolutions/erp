// Package product provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package product

import (
	"context"
	"fmt"
	"time"

	"github.com/deliveranceTechSolutions/erp/business/sys/auth"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Core manages the set of API's for product access.
type Core struct {
	log  *zap.SugaredLogger
	product product.Store
}

// NewCore constructs a core for product api access.
func NewCore(log *zap.SugaredLogger, db *sqlx.DB) Core {
	return Core{
		log:  log,
		product: product.NewStore(log, db),
	}
}

// Create inserts a new product into the database.
func (c Core) Create(ctx context.Context, nu product.Newproduct, now time.Time) (product.product, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.product.Create(ctx, nu, now)
	if err != nil {
		return product.product{}, fmt.Errorf("create: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// Update replaces a product document in the database.
func (c Core) Update(ctx context.Context, claims auth.Claims, productID string, uu product.Updateproduct, now time.Time) error {

	// PERFORM PRE BUSINESS OPERATIONS

	if err := c.product.Update(ctx, claims, productID, uu, now); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}

// Delete removes a product from the database.
func (c Core) Delete(ctx context.Context, claims auth.Claims, productID string) error {

	// PERFORM PRE BUSINESS OPERATIONS

	if err := c.product.Delete(ctx, claims, productID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}

// Query retrieves a list of existing products from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]product.product, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	products, err := c.product.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return products, nil
}

// QueryByID gets the specified product from the database.
func (c Core) QueryByID(ctx context.Context, claims auth.Claims, productID string) (product.product, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.product.QueryByID(ctx, claims, productID)
	if err != nil {
		return product.product{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// QueryByEmail gets the specified product from the database by email.
func (c Core) QueryByEmail(ctx context.Context, claims auth.Claims, email string) (product.product, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.product.QueryByID(ctx, claims, email)
	if err != nil {
		return product.product{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// Authenticate finds a product by their email and verifies their password. On
// success it returns a Claims product representing this product. The claims can be
// used to generate a token for future authentication.
func (c Core) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	claims, err := c.product.Authenticate(ctx, now, email, password)
	if err != nil {
		return auth.Claims{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return claims, nil
}
