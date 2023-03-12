package server

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"auth-service/api/internal/handler"
	"auth-service/api/internal/repository"
	"auth-service/api/internal/service"
)

func Inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	/*
	 * repository layer
	 */
	userRepository := repository.NewUserRepository(d.DB)
	tokenRepository := repository.NewTokenRepository(d.RedisClient)

	/*
	 * repository layer
	 */
	userService := service.NewUserService(&service.USConfig{
		UserRepository: userRepository,
	})

	// load rsa keys
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := os.ReadFile(privKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := os.ReadFile(pubKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// load refresh token secret from env variable
	refreshSecret := os.Getenv("REFRESH_SECRET")

	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse env to int: %v", err)
	}
	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse env to int: %v", err)
	}

	tokenService := service.NewTokenService(&service.TSConfig{
		TokenRepository:       tokenRepository,
		PrivateKey:            privKey,
		PublicKey:             pubKey,
		RefreshSecret:         refreshSecret,
		RefreshExpirationSecs: refreshExp,
		IDExpirationSecs:      idExp,
	})

	// initialize gin.Engine
	router := gin.Default()

	baseURL := os.Getenv("ACCOUNT_API_URL")
	timeoutHandler := os.Getenv("HANDLER_TIMEOUT")

	timeout, err := strconv.ParseInt(timeoutHandler, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse env to int: %v", err)
	}

	handler.NewHandler(&handler.Config{
		R:               router,
		UserService:     userService,
		TokenService:    tokenService,
		BaseURL:         baseURL,
		TimeoutDuration: time.Duration(timeout) * time.Second,
	})

	return router, nil
}
