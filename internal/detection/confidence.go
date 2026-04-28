package detection

type Strategy string

const (
	Email               Strategy = "email"
	Phone               Strategy = "phone"
	Postcode            Strategy = "postcode"
	Name                Strategy = "name"
	FirstName           Strategy = "firstName"
	Surname             Strategy = "surname"
	VehicleRegistration Strategy = "vehicleRegistration"
	Redact              Strategy = "redact"
	Remove              Strategy = "remove"
	None                Strategy = "none"
)

type Confidence int

const (
	NoMatch Confidence = iota
	ValuePatternMatch
	ContainsMatch
	FallbackMatch
	StrongMatch
	ExplicitMatch
)

type Match struct {
	Strategy   Strategy
	Confidence Confidence
	MatchedBy  string
}

func (m Match) Found() bool {
	return m.Strategy != None && m.Confidence != NoMatch
}

func NoMatchResult() Match {
	return Match{Strategy: None, Confidence: NoMatch}
}
