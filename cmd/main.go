package main

import (
	"fmt"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println(cfg)
}
