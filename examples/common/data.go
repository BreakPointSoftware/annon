package common

type Customer struct {
	CustomerID string         `json:"customerId" yaml:"customerId"`
	Email      string         `json:"email" yaml:"email"`
	Phone      string         `json:"phone" yaml:"phone"`
	Postcode   string         `json:"postcode" yaml:"postcode"`
	Vehicle    Vehicle        `json:"vehicle" yaml:"vehicle"`
	Contact    ContactDetails `json:"contact" yaml:"contact"`
	Notes      string         `json:"notes" yaml:"notes" anonymise:"false"`
	Secret     string         `json:"secret" yaml:"secret" anonymise:"remove"`
}

type Vehicle struct {
	Registration string `json:"registration" yaml:"registration"`
}

type ContactDetails struct {
	Name  string `json:"name" yaml:"name"`
	Email string `json:"backupEmail" yaml:"backupEmail"`
}

func DemoCustomer() Customer {
	return Customer{
		CustomerID: "12345",
		Email:      "greg@example.com",
		Phone:      "07700 900123",
		Postcode:   "TN9 1XA",
		Vehicle: Vehicle{
			Registration: "AB12 CDE",
		},
		Contact: ContactDetails{
			Name:  "Greg Bryant",
			Email: "backup@example.com",
		},
		Notes:  "Do not redact this note",
		Secret: "internal-only-secret",
	}
}

func MalformedJSON() []byte {
	return []byte(`{"email":`)
}
