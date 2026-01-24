package initialization

import "log"

func Initialize() {
	if !globalConfigExists() {
		err := createGlobalConfig()
		if err != nil {
			log.Fatal("Failed to create global config: ", err)
		}
	}
}
