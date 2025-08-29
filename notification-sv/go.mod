module github.com/ipavlov93/universe-demo/notification-sv

go 1.24.6

replace github.com/ipavlov93/universe-demo/universe-pkg => ../universe-pkg

replace github.com/ipavlov93/universe-demo/product-eventbus-pkg => ../product-eventbus-pkg

require (
	github.com/aws/aws-sdk-go-v2 v1.38.2
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.2
	github.com/ipavlov93/universe-demo/product-eventbus-pkg v0.0.0-20250828161853-3bc8c30741e3
	github.com/ipavlov93/universe-demo/universe-pkg v0.0.0-20250828161853-3bc8c30741e3
	github.com/stretchr/testify v1.11.1
	go.uber.org/zap v1.27.0
)

require (
	github.com/aws/aws-sdk-go-v2/config v1.31.5 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.18.9 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.29.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.34.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.38.1 // indirect
	github.com/aws/smithy-go v1.23.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
