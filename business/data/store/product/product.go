// Package product contains product related CRUD functionality.
package product

import (
	"context"
	"fmt"
	"time"

	"github.com/deliveranceTechSolutions/erp/business/sys/auth"
	"github.com/deliveranceTechSolutions/erp/business/sys/database"
	"github.com/deliveranceTechSolutions/erp/business/sys/validate"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Store manages the set of API's for product access.
type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

// NewStore constructs a product store for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// Let the client tell you the time, so we can test it correctly.

// Create inserts a new product into the database.
func (s Store) Create(ctx context.Context, np NewProduct, now time.Time) (Product, error) {
	if err := validate.Check(np); err != nil {
		return Product{}, fmt.Errorf("validating data: %w", err)
	}

	prod := Product{
		ID:           validate.GenerateID(),
		Name:         np.Name,
		Quantity:	  np.Quantity,
		OwnerID:	  np.OwnerID,
		MinRetainer:  np.MinRetainer,
		StarterBid:	  np.StarterBid,
		Bids:		  np.Bids,
		DateCreated:  now,
		DateUpdated:  now,
	}

	// If db can return an id then you can use the db to generate the ids
	const q = `
	INSERT INTO products(
		product_id, 
		name, 
		quantity, 
		owner_id, 
		min_retainer, 
		bids, 
		date_created, 
		date_updated
	)VALUES(
		:product_id, 
		:name, 
		:quantity, 
		:owner_id, 
		:min_retainer, 
		:bids, 
		:date_created, 
		:date_updated
	)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, prod); err != nil {
		return Product{}, fmt.Errorf("inserting product: %w", err)
	}

	return prod, nil
}

// Update replaces a product document in the database.
func (s Store) Update(ctx context.Context, claims auth.Claims, productID string, up UpdateProduct, now time.Time) error {
	if err := validate.CheckID(productID); err != nil {
		return database.ErrInvalidID
	}
	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	// we query for the product to verify if the state is full
	// doesn't matter if race conditions win, we need
	// to verify the state is qualified before we UPDATE.
	prod, err := s.QueryByID(ctx, claims, productID)
	if err != nil {
		return fmt.Errorf("updating product productID[%s]: %w", productID, err)
	}

	if up.Name != nil {
		prod.Name = *up.Name
	}
	if up.MinRetainer != nil {
		prod.MinRetainer = *up.MinRetainer
	}
	if up.StarterBid != nil && *up.StarterBid <= *up.MinRetainer {
		prod.MinRetainer = *up.MinRetainer
	}
	
	prod.DateUpdated = now

	const q = `
	UPDATE
		products
	SET 
		"name" = :name,
		"quantity" = :quantity,
		"owner_id" = :owner_id,
		"min_retainer" = :min_retainer,
		"date_updated" = :date_updated
	WHERE
		product_id = :product_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, prod); err != nil {
		return fmt.Errorf("updating productID[%s]: %w", productID, err)
	}

	return nil
}

// Delete removes a product from the database.
func (s Store) Delete(ctx context.Context, claims auth.Claims, productID string) error {
	if err := validate.CheckID(productID); err != nil {
		return database.ErrInvalidID
	}

	// If you are not an admin and looking to delete someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != productID {
		return database.ErrForbidden
	}

	data := struct {
		ProductID string `db:"product_id"`
	}{
		ProductID: productID,
	}

	const q = `
	DELETE FROM
		products
	WHERE
		product_id = :product_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting productID[%s]: %w", productID, err)
	}

	return nil
}

// Query retrieves a list of existing products from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Product, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		products
	ORDER BY
		product_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var products []Product
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &products); err != nil {
		if err == database.ErrNotFound {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("selecting products: %w", err)
	}

	return products, nil
}

// QueryByID gets the specified product from the database.
func (s Store) QueryByID(ctx context.Context, claims auth.Claims, productID string) (Product, error) {
	if err := validate.CheckID(productID); err != nil {
		return Product{}, database.ErrInvalidID
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != productID {
		return Product{}, database.ErrForbidden
	}

	data := struct {
		ProductID string `db:"product_id"`
	}{
		ProductID: productID,
	}

	const q = `
	SELECT
		*
	FROM
		products
	WHERE 
		product_id = :product_id`

	var prod Product
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &prod); err != nil {
		if err == database.ErrNotFound {
			return Product{}, database.ErrNotFound
		}
		return Product{}, fmt.Errorf("selecting productID[%q]: %w", productID, err)
	}

	return prod, nil
}

// QueryMultipleQuantities gets the specified product from the database by email.
func (s Store) QueryMultipleQuantities(ctx context.Context, claims auth.Claims) (Product, error) {
	const q = `
	SELECT
		*
	FROM
		products
	WHERE
		quantity > 1`

	var prod Product
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &prod); err != nil {
		if err == database.ErrNotFound {
			return Product{}, database.ErrNotFound
		}
		return Product{}, fmt.Errorf("selecting email[%q]: %w", email, err)
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != prod.ID {
		return Product{}, database.ErrForbidden
	}

	return prod, nil
}

// QuerySingleQuantities gets the specified product from the database by email.
func (s Store) QuerySingleQuantity(ctx context.Context, claims auth.Claims) (Product, error) {
	const q = `
	SELECT
		*
	FROM
		products
	WHERE
		quantity = 1`

	var prod Product
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &prod); err != nil {
		if err == database.ErrNotFound {
			return Product{}, database.ErrNotFound
		}
		return Product{}, fmt.Errorf("selecting email[%q]: %w", email, err)
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != prod.ID {
		return Product{}, database.ErrForbidden
	}

	return prod, nil
}

// Authenticate finds a product by their email and verifies their password. On
// success it returns a Claims Product representing this product. The claims can be
// used to generate a token for future authentication.
func (s Store) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	const q = `
	SELECT
		*
	FROM
		products
	WHERE
		email = :email`

	var prod Product
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &prod); err != nil {
		if err == database.ErrNotFound {
			return auth.Claims{}, database.ErrNotFound
		}
		return auth.Claims{}, fmt.Errorf("selecting product[%q]: %w", email, err)
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(prod.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, database.ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the product
	// and generate their token.
	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   prod.ID,
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
		Roles: prod.Roles,
	}

	return claims, nil
}


// Dashboard find the product by email and hydrates the view with a config JSON. On
// success it returns a JSON object representing the product's conf. The claims can be
// used to generate a token for future authentication.
func (s Store) Dashboard(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	const q = `
	SELECT
		dash
	FROM
		products
	WHERE
		email = :email`

	var prod Product
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &prod); err != nil {
		if err == database.ErrNotFound {
			return auth.Claims{}, database.ErrNotFound
		}
		return auth.Claims{}, fmt.Errorf("selecting product[%q]: %w", email, err)
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(prod.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, database.ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the product
	// and generate their token.
	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   prod.ID,
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
		Roles: prod.Roles,
	}

	return claims, nil
}