package route

const (
	Angola Country = iota
	SouthArabia
	Turkey
	Russia
	France
	Spain
	China
)

var Countries = map[Country]string{
	Angola:      "Angola",
	SouthArabia: "SouthArabia",
	Turkey:      "Turkey",
	Russia:      "Russia",
	France:      "France",
	Spain:       "Spain",
	China:       "China",
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

type Country int

func (c Country) String() string {
	name, ok := Countries[c]
	if !ok {
		return "Undefined"
	}
	return name
}
