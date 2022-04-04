package dragonSpider

import "os"

//CreateDir creates a directory if it doesn't already exist'
func (ds *DragonSpider) CreateDir(path string) error {
	const mode = 0755

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}
