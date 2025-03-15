package models

type Employee struct {
	ID          int
	FirstName   string
	LastName    string
	Position    string
	Description string
	Image       string
}
type Director struct {
	ID          int
	FirstName   string
	LastName    string
	Position    string
	Description string
	Image       string
}
type Service struct {
	ID          int
	Name        string
	Description string
	Price       int
}

type Review struct {
	ID     int
	Author string
	Text   string
}
type Project struct {
	ID          int
	Name        string
	Description string
	Images      []string
	Videos      []string
}
type Page struct {
	Title     string
	Employees []Employee
	Director  Director
	Services  []Service
	Reviews   []Review
	Projects  []Project
	Employee  Employee
	Name      string
	Email     string
	Message   string
	PrevPage  int
	NextPage  int
}
