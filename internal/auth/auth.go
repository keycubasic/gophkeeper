// Package auth provides authorization.
package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/keycubasic/gophkeeper/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const contextKey = "username"

// JWTManager stores secret key and token duration.
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTManager creates mew JWTManager
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

// Generate generates access token.
func (manager *JWTManager) Generate(user *model.User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// UserClaims stores user claims.
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

// Verify verifies access token.
func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.secretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// AuthInterceptor stores JWTManager.
type AuthInterceptor struct {
	JwtManager *JWTManager
}

// NewAuthInterceptor creates new AuthInterceptor.
func NewAuthInterceptor(jwtManager *JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

// Unary returns UnaryServerInterceptor.
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if strings.Contains(info.FullMethod, "AuthService") {
			return handler(ctx, req)
		}

		newCtx, err := interceptor.authorize(ctx)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

// Stream returns StreamServerInterceptor.
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if strings.Contains(info.FullMethod, "AuthService") {
			return handler(srv, stream)
		}

		newCtx, err := interceptor.authorize(stream.Context())
		if err != nil {
			return err
		}

		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, stream)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := interceptor.JwtManager.Verify(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	ctx = context.WithValue(ctx, contextKey, claims.Username)

	return ctx, nil
}
