package tr

import (
	"log"

	"github.com/xmhrnx/giode/i18n"
)

func init() {
	if err := i18n.Load(files, "locales"); err != nil {
		log.Printf("translations: %v", err)
	}
}
