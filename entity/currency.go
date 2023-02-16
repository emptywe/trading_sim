package entity

type Currency struct {
	Cid   int     `db:"id"`
	Name  string  `db:"name"`
	Value float64 `db:"value"`
}

type CurrencyOutput struct {
	Name  string  `db:"name"`
	Value float64 `db:"value"`
}
