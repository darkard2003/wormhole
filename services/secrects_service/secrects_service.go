package secrectsservice

import "github.com/zalando/go-keyring"

type SecretsService struct {
	service string
}

func NewSecretsService() *SecretsService {
	return &SecretsService{
		service: "wormhole-cli",
	}
}

func (s *SecretsService) SetSecret(key, value string) error {
	return keyring.Set(s.service, key, value)
}

func (s *SecretsService) GetSecret(key string) (string, error) {
	return keyring.Get(s.service, key)
}

func (s *SecretsService) DeleteSecret(key string) error {
	return keyring.Delete(s.service, key)
}

func (s *SecretsService) GetJWTSecret() (string, error) {
	secrect, err := s.GetSecret("JWT_SECRET_KEY")
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}
		return "", err
	}
	return secrect, nil
}

func (s *SecretsService) SetJWTSecret(secret string) error {
	return s.SetSecret("JWT_SECRET_KEY", secret)
}

func (s *SecretsService) GetRefreshTokenSecret() (string, error) {
	secrect, err := s.GetSecret("REFRESH_TOKEN_SECRET_KEY")
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}
		return "", err
	}

	return secrect, nil
}

func (s *SecretsService) SetRefreshTokenSecret(secret string) error {
	return s.SetSecret("REFRESH_TOKEN_SECRET_KEY", secret)
}
