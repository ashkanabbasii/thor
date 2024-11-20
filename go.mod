module github.com/ashkanabbasii/thor

go 1.22

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1
	github.com/ethereum/go-ethereum v1.14.12
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.1
	github.com/hashicorp/golang-lru v0.0.0-20160813221303-0a025b7e63ad
	github.com/holiman/uint256 v1.3.1
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.18.0
	github.com/prometheus/client_model v0.5.0
	github.com/stretchr/testify v1.9.0
	github.com/syndtr/goleveldb v1.0.1-0.20220614013038-64ee5596c38a
	github.com/vechain/go-ecvrf v0.0.0-20220525125849-96fa0442e765
	github.com/vechain/thor/v2 v2.1.4
	golang.org/x/crypto v0.22.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/qianbin/directcache v0.9.7 // indirect
	golang.org/x/sys v0.22.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/syndtr/goleveldb => github.com/vechain/goleveldb v1.0.1-0.20220809091043-51eb019c8655

//replace github.com/ethereum/go-ethereum => github.com/vechain/go-ethereum v1.8.15-0.20240528020007-2994c2a24b9c
