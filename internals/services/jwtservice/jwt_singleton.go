package jwtservice

import "log"

var jwtServiceInstance *JWTService

func GetJWTService() *JWTService {
	if jwtServiceInstance == nil {
		log.Fatal("jwtServiceInstance is nil")
	}
	return jwtServiceInstance
}

func InitJWTService(accessSecrectKey, refreshSecrectKey string) {
	jwtServiceInstance = NewJWTService(accessSecrectKey, refreshSecrectKey)
}
