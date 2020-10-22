module github.com/byuoitav/av-control-api

go 1.14

require (
	github.com/byuoitav/adcp-driver v0.1.1
	github.com/byuoitav/atlona-driver v1.5.7
	github.com/byuoitav/av-control-api/api v0.0.0-20200824162301-775202bed269
	github.com/byuoitav/justaddpower-driver v0.1.3
	github.com/byuoitav/kramer v0.0.0-00010101000000-000000000000
	github.com/byuoitav/kramer-driver v0.1.12
	github.com/byuoitav/qsc-driver v0.1.6
	github.com/byuoitav/sonyrest-driver v0.1.8
	github.com/byuoitav/wspool v0.1.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-kivik/couchdb/v3 v3.1.0
	github.com/go-kivik/kivik/v3 v3.1.1
	github.com/go-kivik/kivikmock/v3 v3.1.1
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/google/go-cmp v0.5.0
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/matryer/is v1.4.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/segmentio/ksuid v1.0.2
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20200204104054-c9f3fb736b72 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/byuoitav/kramer => ../kramer
