//go:build !k8s

package config

var Config = blueBookConfig{
	DB{DSN: "root:root@tcp(localhost:13326)/bluebook"},
}
