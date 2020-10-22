module github.com/byuoitav/av-control-api

go 1.14

require (
	github.com/byuoitav/adcp-driver v0.1.1
	github.com/byuoitav/atlona-driver v1.5.8
	github.com/byuoitav/av-control-api/api v0.0.0-20200824162301-775202bed269
	github.com/byuoitav/common v0.0.0-20200521193927-1fdf4e0a4271 // indirect
	github.com/byuoitav/justaddpower-driver v0.1.3
	github.com/byuoitav/keydigital-driver v0.0.10
	github.com/byuoitav/kramer v0.0.0-00010101000000-000000000000
	github.com/byuoitav/kramer-driver v0.1.13
	github.com/byuoitav/london-driver v0.1.3
	github.com/byuoitav/qsc-driver v0.1.7
	github.com/byuoitav/sonyrest-driver v0.1.8
	github.com/byuoitav/wspool v0.1.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-kivik/couchdb/v3 v3.2.1
	github.com/go-kivik/kivik/v3 v3.2.0
	github.com/go-kivik/kivikmock/v3 v3.1.1
	github.com/go-playground/validator/v10 v10.4.0 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/go-cmp v0.5.0
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/segmentio/ksuid v1.0.3
	github.com/spf13/pflag v1.0.5
	github.com/ugorji/go v1.1.12 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb
	golang.org/x/sync v0.0.0-20201008141435-b3e1573b7520
	golang.org/x/sys v0.0.0-20201015000850-e3ed0017c211 // indirect
	golang.org/x/tools v0.0.0-20201015182029-a5d9e455e9c4 // indirect
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
	honnef.co/go/tools v0.0.1-2020.1.6 // indirect
)

replace github.com/byuoitav/kramer => ../kramer
