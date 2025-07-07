package dbservice

var dbServiceInstance *DBService

func GetDBService() *DBService {
	if dbServiceInstance == nil {
		dbServiceInstance = &DBService{
			DB:          nil,
			Initialized: false,
		}
	}
	return dbServiceInstance
}