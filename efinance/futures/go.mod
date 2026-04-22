module github.com/efinance/efinance/efinance/futures

go 1.21

require (
	github.com/efinance/efinance/efinance/common v0.0.0
	github.com/efinance/efinance/efinance/errors v0.0.0
)

replace (
	github.com/efinance/efinance/efinance/common => ../common
	github.com/efinance/efinance/efinance/errors => ../errors
)
