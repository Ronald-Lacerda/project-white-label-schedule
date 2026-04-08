package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"schedule/internal/shared"
)

type Claims struct {
	UserID          string `json:"user_id"`
	EstablishmentID string `json:"establishment_id"`
	jwt.RegisteredClaims
}

type refreshPayload struct {
	UserID          string `json:"user_id"`
	EstablishmentID string `json:"establishment_id"`
}

type Service struct {
	repo      Repository
	redis     *redis.Client
	jwtSecret []byte
	jwtExpiry time.Duration
	rtExpiry  time.Duration
}

func NewService(repo Repository, rdb *redis.Client, jwtSecret string, jwtExpiry, rtExpiry time.Duration) *Service {
	return &Service{
		repo:      repo,
		redis:     rdb,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: jwtExpiry,
		rtExpiry:  rtExpiry,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, shared.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, shared.ErrUnauthorized
	}

	return s.issueTokenPair(ctx, user)
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*TokenPair, error) {
	ownerName := strings.TrimSpace(input.OwnerName)
	establishmentName := strings.TrimSpace(input.EstablishmentName)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := input.Password
	slug := sanitizeSlug(input.Slug)

	if ownerName == "" || establishmentName == "" || email == "" || password == "" || slug == "" {
		return nil, shared.ErrInvalidInput
	}
	if len(password) < 8 || !isValidEmail(email) {
		return nil, shared.ErrInvalidInput
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var contactPhone *string
	if input.ContactPhone != nil {
		trimmedPhone := strings.TrimSpace(*input.ContactPhone)
		if trimmedPhone != "" {
			contactPhone = &trimmedPhone
		}
	}

	user, err := s.repo.CreateAccount(ctx, CreateAccountInput{
		OwnerName:         ownerName,
		EstablishmentName: establishmentName,
		Email:             email,
		PasswordHash:      string(passwordHash),
		Slug:              slug,
		Timezone:          "America/Sao_Paulo",
		ContactPhone:      contactPhone,
	})
	if err != nil {
		return nil, err
	}

	return s.issueTokenPair(ctx, user)
}

func (s *Service) Logout(ctx context.Context, rawRefresh string) error {
	key := s.refreshKey(rawRefresh)
	return s.redis.Del(ctx, key).Err()
}

func (s *Service) Refresh(ctx context.Context, rawRefresh string) (*TokenPair, error) {
	key := s.refreshKey(rawRefresh)

	val, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, shared.ErrUnauthorized
	}

	var payload refreshPayload
	if err := json.Unmarshal([]byte(val), &payload); err != nil {
		return nil, shared.ErrUnauthorized
	}

	user, err := s.repo.FindUserByID(ctx, payload.UserID)
	if err != nil {
		return nil, shared.ErrUnauthorized
	}

	s.redis.Del(ctx, key)
	return s.issueTokenPair(ctx, user)
}

func (s *Service) ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, shared.ErrUnauthorized
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, shared.ErrUnauthorized
	}
	return claims, nil
}

func (s *Service) issueTokenPair(ctx context.Context, user *User) (*TokenPair, error) {
	accessToken, err := s.createAccessToken(user)
	if err != nil {
		return nil, err
	}

	rawRefresh, err := generateToken(32)
	if err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(refreshPayload{
		UserID:          user.ID,
		EstablishmentID: user.EstablishmentID,
	})

	key := s.refreshKey(rawRefresh)
	if err := s.redis.Set(ctx, key, payload, s.rtExpiry).Err(); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		User:         user,
	}, nil
}

func (s *Service) createAccessToken(user *User) (string, error) {
	claims := Claims{
		UserID:          user.ID,
		EstablishmentID: user.EstablishmentID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) refreshKey(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return fmt.Sprintf("refresh:%x", hash)
}

func generateToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func isValidEmail(email string) bool {
	return emailPattern.MatchString(email)
}

func sanitizeSlug(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	var builder strings.Builder
	lastWasHyphen := false

	for _, r := range input {
		if replacement, ok := slugAccentMap[r]; ok {
			r = replacement
		}

		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
			lastWasHyphen = false
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
			lastWasHyphen = false
		case unicode.IsSpace(r) || r == '-' || r == '_' || r == '/':
			if builder.Len() > 0 && !lastWasHyphen {
				builder.WriteByte('-')
				lastWasHyphen = true
			}
		}
	}

	return strings.Trim(builder.String(), "-")
}

var slugAccentMap = map[rune]rune{
	'á': 'a',
	'à': 'a',
	'â': 'a',
	'ã': 'a',
	'ä': 'a',
	'é': 'e',
	'è': 'e',
	'ê': 'e',
	'ë': 'e',
	'í': 'i',
	'ì': 'i',
	'î': 'i',
	'ï': 'i',
	'ó': 'o',
	'ò': 'o',
	'ô': 'o',
	'õ': 'o',
	'ö': 'o',
	'ú': 'u',
	'ù': 'u',
	'û': 'u',
	'ü': 'u',
	'ç': 'c',
	'ñ': 'n',
}
