package revoluchat

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// JWTManager handles token generation and JWKS serving.
type JWTManager struct {
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
	kid       string
	serverKey string // Secure key for JWKS endpoint validation
}

// NewJWTManager creates a new JWTManager and computes the KID thumbprint.
func NewJWTManager(privKeyPath, pubKeyPath, serverKey string) (*JWTManager, error) {
	signBytes, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	verifyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Compute Thumbprint KID
	key, _ := jwk.FromRaw(verifyKey)
	tp, _ := key.Thumbprint(crypto.SHA256)
	kid := base64.RawURLEncoding.EncodeToString(tp)

	return &JWTManager{
		signKey:   signKey,
		verifyKey: verifyKey,
		kid:       kid,
		serverKey: serverKey,
	}, nil
}

// GenerateToken creates a signed JWT with KID in header.
func (m *JWTManager) GenerateToken(userID string, appID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":    userID,
		"app_id": appID,
		"iss":    "revolu-be",
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = m.kid
	return token.SignedString(m.signKey)
}

// JWKSHandler returns a handler function for JWKS endpoint with security check.
func (m *JWTManager) JWKSHandler(w http.ResponseWriter, r *http.Request) {
	// Security Check: Validate X-Server-Key header
	clientKey := r.Header.Get("X-Server-Key")
	if clientKey == "" {
		clientKey = r.URL.Query().Get("server_key")
	}
	if m.serverKey != "" && clientKey != m.serverKey {
		http.Error(w, "Unauthorized: Invalid server key", http.StatusUnauthorized)
		return
	}

	key, err := jwk.FromRaw(m.verifyKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = key.Set(jwk.KeyIDKey, m.kid)
	_ = key.Set(jwk.AlgorithmKey, "RS256")
	_ = key.Set(jwk.KeyUsageKey, "sig")

	jwks := map[string]interface{}{
		"keys": []interface{}{key},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}
