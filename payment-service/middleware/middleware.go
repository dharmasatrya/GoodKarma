package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var protectedMethods = map[string]bool{
	"/payment.PaymentService/CreateWallet":          false,
	"/payment.PaymentService/UpdateWalletBalance":   true,
	"/payment.PaymentService/Withdraw":              true,
	"/payment.PaymentService/CreateInvoice":         true,
	"/payment.PaymentService/GetWalletByUserId":     true,
	"/payment.PaymentService/XenditInvoiceCallback": false,
}

func UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Check if the method requires authentication
	if requiresAuth := protectedMethods[info.FullMethod]; !requiresAuth {
		// Method doesn't require auth, proceed without checking
		return handler(ctx, req)
	}

	// Method requires auth, proceed with authentication
	ctx, err := AuthInterceptor(ctx)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func AuthInterceptor(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("No metadata found")
		return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	log.Printf("Metadata received: %v", md)

	token := md["authorization"]
	if len(token) == 0 {
		log.Println("Missing token")
		return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	// Remove this line that checks for exact match
	// if token[0] != "Bearer valid-token" {  <- This was the problem

	// Instead, validate the JWT token
	tokenString := strings.TrimPrefix(token[0], "Bearer ")
	claims, err := validateToken(tokenString)
	if err != nil {
		log.Println("Invalid token:", err)
		return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	// Optional: Add claims to context if needed
	newCtx := context.WithValue(ctx, "claims", claims)

	log.Println("Token validated successfully")
	return newCtx, nil
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Use the same secret as your auth service
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil // Replace with your actual secret
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
