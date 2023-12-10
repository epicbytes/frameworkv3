package vault

import (
	vault "github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"os"
)

func New() {
	config := vault.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")
	//fmt.Println(config)
	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatal().AnErr("Unable to initialize a Vault client: ", err).Send()
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	//kvv := client.KVv2("configuration")

	//fmt.Println(kvv)

}
