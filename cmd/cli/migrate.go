package main

func runMigration(migrationType, step string) error {
	dsn := GetDSN()

	switch migrationType {
	case "up":
		err := ds.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "down":
		if step == "all" {
			err := ds.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := ds.Steps(dsn, -1)
			if err != nil {
				return err
			}
		}

	case "reset":
		err := ds.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = ds.MigrateUp(dsn)
		if err != nil {
			return err
		}

	default:
		showHelp()
	}

	return nil
}
