module github.com/byuoitav/av-control-api/api

go 1.14

require (
	github.com/byuoitav/av-control-api v0.3.9
	github.com/gin-gonic/gin v1.6.3
	github.com/go-kivik/couchdb/v3 v3.1.0
	github.com/go-kivik/kivik/v3 v3.1.1
	github.com/go-kivik/kivikmock/v3 v3.1.1
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/goccy/go-graphviz v0.0.5
	github.com/google/go-cmp v0.5.0
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/segmentio/ksuid v1.0.2
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.14.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	gonum.org/v1/gonum v0.7.0
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	istio.io/pkg v0.0.0-20200630182444-e8a83c9625a3
)

replace github.com/byuoitav/av-control-api => ../
