package cmd

type MigratorCommand struct {
	Migrate MigrateCommand `command:"migrate" alias:"m" description:"Migrate credentials from vars store to CredHub"`
}

var Migrator MigratorCommand
