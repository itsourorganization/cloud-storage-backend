package main

import (
	"github.com/undefeel/cloud-storage-backend/internal/config"
)

func main() {
	//read config
	_ = config.MustLoad()
}
