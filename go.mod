module github.com/c4ei/c4exapi

go 1.14

require (
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/go-pg/pg/v9 v9.1.3
	github.com/golang-migrate/migrate/v4 v4.7.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/jessevdk/go-flags v1.4.0
	github.com/kaspanet/go-secp256k1 v0.0.2
	github.com/c4ei/c4exd v0.6.2
	github.com/pkg/errors v0.9.1
)

replace github.com/c4ei/c4exd => ../kaspad
