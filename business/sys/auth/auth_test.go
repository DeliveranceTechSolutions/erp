package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/deliveranceTechSolutions/erp/business/sys/auth"
	"github.com/golang-jwt/jwt/v4"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestAuth(t *testing.T) {
	// Mission Statment - The What of the Test
	t.Log("Given the need to be able to authenticate and authorize access.")
	{
		testID := 0
		// Frequency - The When of the Test - could have multiple whens next to each other
		t.Logf("\tTest %d:\tWhen handling a single user.", testID)
		{
			// Creat Private Key with KID
			const keyID = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				// If it's a mission critical error the use t.Fatalf() - if not then use t.Errorf()
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a private key: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create a private key.", success, testID)

			// Instantiate New KeyStore to store private key into
			a, err := auth.New(keyID, &keyStore{pk: privateKey})
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create an authenticator: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create an authenticator.", success, testID)

			// Generate claims using the KID and Role
			claims := auth.Claims{
				StandardClaims: jwt.StandardClaims{
					Issuer:    "service project",
					Subject:   "5cf37266-3473-4006-984f-9325122678b7",
					ExpiresAt: time.Now().Add(time.Hour).Unix(),
					IssuedAt:  time.Now().UTC().Unix(),
				},
				Roles: []string{auth.RoleAdmin},
			}

			// Generate Token to so we can parse claims and validate the token
			token, err := a.GenerateToken(claims)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to generate a JWT: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to generate a JWT.", success, testID)

			// Validate Token into parsedClaims 
			parsedClaims, err := a.ValidateToken(token)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to parse the claims: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to parse the claims.", success, testID)

			// analyze the results of the parsed claims using expected/got syntax
			if exp, got := len(claims.Roles), len(parsedClaims.Roles); exp != got {
				t.Logf("\t\tTest %d:\texp: %d", testID, exp)
				t.Logf("\t\tTest %d:\tgot: %d", testID, got)
				t.Fatalf("\t%s\tTest %d:\tShould have the expected number of roles: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected number of roles.", success, testID)

			if exp, got := claims.Roles[0], parsedClaims.Roles[0]; exp != got {
				t.Logf("\t\tTest %d:\texp: %v", testID, exp)
				t.Logf("\t\tTest %d:\tgot: %v", testID, got)
				t.Fatalf("\t%s\tTest %d:\tShould have the expected roles: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected roles.", success, testID)
		}
	}
}

// =============================================================================

type keyStore struct {
	pk *rsa.PrivateKey
}

func (ks *keyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	return ks.pk, nil
}

func (ks *keyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	return &ks.pk.PublicKey, nil
}