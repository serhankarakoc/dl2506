package main

import (
	"flag"

	"davet.link/configs/databaseconfig"
	"davet.link/configs/logconfig"
	"davet.link/database"
)

func main() {
	logconfig.InitLogger()
	defer logconfig.SyncLogger()
	migrateFlag := flag.Bool("migrate", false, "Veritabanı başlatma işlemini çalıştır (migrasyonları içerir)")
	seedFlag := flag.Bool("seed", false, "Veritabanı başlatma işlemini çalıştır (seederları içerir)")
	flag.Parse()

	databaseconfig.InitDB()
	defer databaseconfig.CloseDB()

	db := databaseconfig.GetDB()

	logconfig.SLog.Info("Veritabanı başlatma işlemi çalıştırılıyor...")
	database.Initialize(db, *migrateFlag, *seedFlag)

	logconfig.SLog.Info("Veritabanı başlatma işlemi tamamlandı.")
}
