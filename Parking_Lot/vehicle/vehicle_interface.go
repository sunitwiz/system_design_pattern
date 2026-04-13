package vehicle

type Vehicle interface {
	GetType() VehicleType
	GetLicensePlate() string
	String() string
}
