package services

type Connectable interface {
	Name() string
	CheckConnection() error
}

type ServiceConnectionError struct {
	ServiceName string
	Error       string
}

func CheckAllConnections(services ...Connectable) []ServiceConnectionError {
	var allErrs []ServiceConnectionError
	for _, service := range services {
		if err := service.CheckConnection(); err != nil {
			allErrs = append(allErrs, ServiceConnectionError{
				ServiceName: service.Name(),
				Error:       err.Error(),
			})
		}
	}
	return allErrs
}
