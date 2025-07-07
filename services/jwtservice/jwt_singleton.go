package jwtservice

var jwtServiceInstance *JWTService

func GetJWTService() *JWTService {
	if jwtServiceInstance == nil {
		jwtServiceInstance = NewJWTService()
	}
	return jwtServiceInstance
}
