package route

// Const block defines a set of constants representing different
// countries such as Angola, South Arabia, Turkey, Russia, France, Spain, and China.
const (
	Angola Country = iota
	SouthArabia
	Turkey
	Russia
	France
	Spain
	China
)

// Countries variable is a map that maps each Country constant to its corresponding string name.
var Countries = map[Country]string{
	Angola:      "Angola",
	SouthArabia: "SouthArabia",
	Turkey:      "Turkey",
	Russia:      "Russia",
	France:      "France",
	Spain:       "Spain",
	China:       "China",
}

// Country represents country type.
type Country int

// String defined in the Country type, which returns the string name of the country based on its constant value.
func (c Country) String() string {
	name, ok := Countries[c]
	if !ok {
		return "Undefined"
	}
	return name
}

var countries = []Country{
	Angola,
	SouthArabia,
	Turkey,
	Russia,
	France,
	Spain,
	China,
}
