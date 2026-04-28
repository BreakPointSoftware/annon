package detection

type MatchType string

const (
	Strong   MatchType = "strong"
	Fallback MatchType = "fallback"
	Contains MatchType = "contains"
)

type Rule struct {
	Strategy Strategy
	Type     MatchType
	Fields   []string
	Exclude  []string
}

func StrongRule(strategy Strategy, fields ...string) Rule {
	return Rule{Strategy: strategy, Type: Strong, Fields: fields}
}

func FallbackRule(strategy Strategy, fields ...string) Rule {
	return Rule{Strategy: strategy, Type: Fallback, Fields: fields}
}

func ContainsRule(strategy Strategy, fields []string, exclude []string) Rule {
	return Rule{Strategy: strategy, Type: Contains, Fields: fields, Exclude: exclude}
}

func DefaultRules() []Rule {
	return []Rule{
		StrongRule(Email, "email", "emailAddress"),
		StrongRule(Phone, "phoneNumber", "mobileNumber", "telephoneNumber"),
		StrongRule(Postcode, "postcode", "postCode", "postalCode"),
		StrongRule(FirstName, "firstName", "forename", "givenName"),
		StrongRule(Surname, "surname", "lastName", "familyName"),
		StrongRule(VehicleRegistration, "vehicleRegistration", "vehicleReg", "vrm"),
		FallbackRule(VehicleRegistration, "registration", "reg"),
		FallbackRule(Phone, "phone", "mobile", "telephone"),
		ContainsRule(Name, []string{"name"}, []string{"username", "filename", "hostname", "domainname"}),
		ContainsRule(Email, []string{"email"}, nil),
	}
}
