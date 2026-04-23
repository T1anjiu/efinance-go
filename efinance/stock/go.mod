module github.com/T1anjiu/efinance-go/efinance/stock

go 1.21

require (
	github.com/T1anjiu/efinance-go/efinance/common v0.0.0
	github.com/T1anjiu/efinance-go/efinance/errors v0.0.0
)

replace (
	github.com/T1anjiu/efinance-go/efinance/common => ../common
	github.com/T1anjiu/efinance-go/efinance/errors => ../errors
)
