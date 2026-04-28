package strategy

type Strategy interface {
	Name() string
	Anonymise(value any, ctx Context) (any, error)
}

func DefaultStrategies() []Strategy {
	return []Strategy{
		EmailStrategy{},
		PhoneStrategy{},
		PostcodeStrategy{},
		NameStrategy{strategyName: "name"},
		NameStrategy{strategyName: "firstName"},
		NameStrategy{strategyName: "surname"},
		VehicleStrategy{},
		RedactStrategy{},
	}
}
