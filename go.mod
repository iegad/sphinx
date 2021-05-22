module github.com/iegad/sphinx

go 1.16

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.2.0
	github.com/iegad/hydra v0.0.1
	github.com/iegad/kraken v0.0.1
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

replace (
	github.com/iegad/hydra v0.0.1 => ../hydra
	github.com/iegad/kraken v0.0.1 => ../kraken
)
