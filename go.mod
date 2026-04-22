module github.com/T1anjiu/efinance-go

go 1.21

require (
	github.com/T1anjiu/efinance-go/efinance/bond v0.0.0-00010101000000-000000000000
	github.com/T1anjiu/efinance-go/efinance/common v0.0.0
	github.com/T1anjiu/efinance-go/efinance/fund v0.0.0-00010101000000-000000000000
	github.com/T1anjiu/efinance-go/efinance/futures v0.0.0-00010101000000-000000000000
	github.com/T1anjiu/efinance-go/efinance/stock v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
)

require github.com/T1anjiu/efinance-go/efinance/errors v0.0.0 // indirect

replace (
	github.com/T1anjiu/efinance-go/efinance/bond => ./efinance/bond
	github.com/T1anjiu/efinance-go/efinance/common => ./efinance/common
	github.com/T1anjiu/efinance-go/efinance/errors => ./efinance/errors
	github.com/T1anjiu/efinance-go/efinance/fund => ./efinance/fund
	github.com/T1anjiu/efinance-go/efinance/futures => ./efinance/futures
	github.com/T1anjiu/efinance-go/efinance/stock => ./efinance/stock
)