module sqs-fargate-consumer

go 1.22.0

toolchain go1.22.9

require (
	github.com/aws/aws-sdk-go-v2 v1.32.4
	github.com/aws/aws-sdk-go-v2/config v1.28.3
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.42.4
	github.com/aws/aws-sdk-go-v2/service/sqs v1.37.0
	github.com/gdcorp-domains/fulfillment-go-api v1.0.88
	github.com/gdcorp-domains/fulfillment-gosecrets v1.0.16
	github.com/google/uuid v1.6.0
)

require (
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-lambda-go v1.36.0 // indirect
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.44 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.32.4 // indirect
	github.com/aws/aws-xray-sdk-go v1.1.0 // indirect
	github.com/aws/smithy-go v1.22.0 // indirect
	github.com/awslabs/aws-lambda-go-api-proxy v0.13.3 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874 // indirect
	github.com/elastic/go-licenser v0.3.1 // indirect
	github.com/elastic/go-sysinfo v1.11.1 // indirect
	github.com/elastic/go-windows v1.0.1 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/gdcorp-domains/fulfillment-gin-writer v1.0.25 // indirect
	github.com/gdcorp-domains/fulfillment-ginutils v1.0.16 // indirect
	github.com/gdcorp-domains/fulfillment-go-filebuffer v1.0.0 // indirect
	github.com/gdcorp-domains/fulfillment-goapimodels v1.0.55 // indirect
	github.com/gdcorp-domains/fulfillment-golang-httpclient v1.0.84 // indirect
	github.com/gdcorp-domains/fulfillment-golang-logging v1.0.21 // indirect
	github.com/gdcorp-domains/fulfillment-golang-middleware v1.0.54 // indirect
	github.com/gdcorp-domains/fulfillment-golang-sso-auth v1.0.60 // indirect
	github.com/gdcorp-domains/fulfillment-gotranslate v1.0.3 // indirect
	github.com/gdcorp-domains/fulfillment-govalidate v1.0.60 // indirect
	github.com/gdcorp-domains/fulfillment-openapihandler v1.0.7 // indirect
	github.com/gdcorp-domains/fulfillment-structiterator v1.0.11 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-contrib/static v0.0.1 // indirect
	github.com/gin-gonic/gin v1.8.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/harlow/kinesis-consumer v0.3.3 // indirect
	github.com/jcchavezs/porto v0.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/juju/xml v0.0.0-20160224194805-b5bf18ebd8b8 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/mr-tron/base58 v1.1.2 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rs/cors v1.8.1 // indirect
	github.com/rs/cors/wrapper/gin v0.0.0-20221003140808-fcebdb403f4d // indirect
	github.com/santhosh-tekuri/jsonschema v1.2.4 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	go.elastic.co/apm v1.15.0 // indirect
	go.elastic.co/apm/module/apmgin v1.11.0 // indirect
	go.elastic.co/apm/module/apmhttp v1.11.0 // indirect
	go.elastic.co/fastjson v1.3.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/exp v0.0.0-20230810033253-352e893a4cad // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	golang.org/x/tools v0.25.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	howett.net/plist v1.0.0 // indirect
)
