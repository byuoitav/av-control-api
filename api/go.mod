module github.com/byuoitav/av-control-api/api

go 1.14

require (
	github.com/byuoitav/av-control-api v0.3.9
	github.com/go-kivik/couchdb v2.0.0+incompatible // indirect
	github.com/go-kivik/couchdb/v3 v3.1.0
	github.com/go-kivik/couchdb/v4 v4.0.0-20200502105845-f8d1cc2b7e9f
	github.com/go-kivik/kivik/v3 v3.1.1
	github.com/go-kivik/kivik/v4 v4.0.0-20200502210153-a9e688f1b1cd
	github.com/go-kivik/kivikmock v2.0.0+incompatible
	github.com/go-kivik/kivikmock/v3 v3.1.1
	github.com/go-kivik/kiviktest v2.0.0+incompatible // indirect
	github.com/goccy/go-graphviz v0.0.5
	github.com/google/go-cmp v0.4.1
	github.com/labstack/echo v3.3.10+incompatible
	github.com/segmentio/ksuid v1.0.2
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.14.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	gonum.org/v1/gonum v0.7.0
	google.golang.org/grpc v1.30.0
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace github.com/byuoitav/av-control-api => ../
