module github.com/efinance/efinance

go 1.21

require (
	github.com/efinance/efinance/efinance/bond v0.0.0-00010101000000-000000000000
	github.com/efinance/efinance/efinance/common v0.0.0
	github.com/efinance/efinance/efinance/fund v0.0.0-00010101000000-000000000000
	github.com/efinance/efinance/efinance/futures v0.0.0-00010101000000-000000000000
	github.com/efinance/efinance/efinance/stock v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
)

require github.com/efinance/efinance/efinance/errors v0.0.0 // indirect

replace (
	github.com/efinance/efinance/efinance/bond => ./efinance/bond
	github.com/efinance/efinance/efinance/common => ./efinance/common
	github.com/efinance/efinance/efinance/errors => ./efinance/errors
	github.com/efinance/efinance/efinance/fund => ./efinance/fund
	github.com/efinance/efinance/efinance/futures => ./efinance/futures
	github.com/efinance/efinance/efinance/stock => ./efinance/stock
)
