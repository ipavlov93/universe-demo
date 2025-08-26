module github.com/ipavlov93/universe-demo/product-service

go 1.24.6

replace github.com/ipavlov93/universe-demo/universe-pkg => ../universe-pkg

require (
	github.com/ipavlov93/universe-demo/universe-pkg v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.0
)

require go.uber.org/multierr v1.10.0 // indirect
