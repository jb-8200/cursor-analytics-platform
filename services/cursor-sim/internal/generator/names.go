package generator

import (
	"fmt"
	"math/rand"
	"strings"
)

// NameGenerator generates deterministic realistic names
type NameGenerator struct {
	rng *rand.Rand
}

// NewNameGenerator creates a new name generator with the given seed
func NewNameGenerator(seed int64) *NameGenerator {
	return &NameGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Common first names for variety
var firstNames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"William", "Barbara", "David", "Elizabeth", "Richard", "Susan", "Joseph", "Jessica",
	"Thomas", "Sarah", "Christopher", "Karen", "Charles", "Lisa", "Daniel", "Nancy",
	"Matthew", "Betty", "Anthony", "Helen", "Mark", "Sandra", "Donald", "Donna",
	"Steven", "Carol", "Andrew", "Ruth", "Paul", "Sharon", "Joshua", "Michelle",
	"Kenneth", "Laura", "Kevin", "Sarah", "Brian", "Kimberly", "George", "Deborah",
	"Timothy", "Jessica", "Ronald", "Shirley", "Edward", "Cynthia", "Jason", "Angela",
	"Jeffrey", "Melissa", "Ryan", "Brenda", "Jacob", "Amy", "Gary", "Anna",
	"Nicholas", "Rebecca", "Eric", "Virginia", "Jonathan", "Kathleen", "Stephen", "Pamela",
	"Larry", "Martha", "Justin", "Debra", "Scott", "Amanda", "Brandon", "Stephanie",
	"Benjamin", "Carolyn", "Samuel", "Christine", "Raymond", "Marie", "Gregory", "Janet",
	"Alexander", "Catherine", "Patrick", "Frances", "Frank", "Ann", "Jack", "Joyce",
	"Dennis", "Diane", "Jerry", "Alice", "Tyler", "Julie", "Aaron", "Heather",
	"Jose", "Teresa", "Adam", "Doris", "Nathan", "Gloria", "Henry", "Evelyn",
	"Zachary", "Jean", "Douglas", "Cheryl", "Peter", "Mildred", "Kyle", "Katherine",
}

// Common last names for variety
var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas",
	"Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White",
	"Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker", "Young",
	"Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores",
	"Green", "Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell",
	"Carter", "Roberts", "Gomez", "Phillips", "Evans", "Turner", "Diaz", "Parker",
	"Cruz", "Edwards", "Collins", "Reyes", "Stewart", "Morris", "Morales", "Murphy",
	"Cook", "Rogers", "Gutierrez", "Ortiz", "Morgan", "Cooper", "Peterson", "Bailey",
	"Reed", "Kelly", "Howard", "Ramos", "Kim", "Cox", "Ward", "Richardson",
	"Watson", "Brooks", "Chavez", "Wood", "James", "Bennett", "Gray", "Mendoza",
	"Ruiz", "Hughes", "Price", "Alvarez", "Castillo", "Sanders", "Patel", "Myers",
	"Long", "Ross", "Foster", "Jimenez", "Powell", "Jenkins", "Perry", "Russell",
	"Sullivan", "Bell", "Coleman", "Butler", "Henderson", "Barnes", "Gonzales", "Fisher",
	"Vasquez", "Simmons", "Romero", "Jordan", "Patterson", "Alexander", "Hamilton", "Graham",
}

// GenerateName generates a deterministic name based on the internal RNG state
func (n *NameGenerator) GenerateName() (firstName, lastName string) {
	firstName = firstNames[n.rng.Intn(len(firstNames))]
	lastName = lastNames[n.rng.Intn(len(lastNames))]
	return firstName, lastName
}

// GenerateEmail generates an email from a first and last name
func (n *NameGenerator) GenerateEmail(firstName, lastName string) string {
	// Convert to lowercase for email
	return fmt.Sprintf("%s.%s@company.com",
		strings.ToLower(firstName),
		strings.ToLower(lastName))
}
