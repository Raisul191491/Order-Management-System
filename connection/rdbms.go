package connection

import (
	"context"
	"fmt"
	"log"
	"oms/config"
	migrations "oms/migration"
	"oms/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	sqlLogger "gorm.io/gorm/logger"
)

func InitDB(cfg config.Config) (*gorm.DB, *gorm.DB) {
	readDataSourceName := GetReadDSN(cfg)
	writeDataSourceName := GetWriteDSN(cfg)

	// Connect to the master database
	masterDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  writeDataSourceName,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: sqlLogger.Default,
	})
	if err != nil {
		log.Fatalf("error initializing master DB instance: %v", err)
	}

	// Connect to the read replica database
	replicaDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  readDataSourceName,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: sqlLogger.Default,
	})
	if err != nil {
		log.Fatalf("error initializing replica DB instance: %v", err)
	}

	// Get underlying sql.DB instances
	sqlMasterDB, err := masterDB.DB()
	if err != nil {
		log.Fatalf("error getting master sql.DB: %v", err)
	}

	sqlReplicaDB, err := replicaDB.DB()
	if err != nil {
		log.Fatalf("error getting replica sql.DB: %v", err)
	}

	// Set connection pool settings
	sqlMasterDB.SetMaxOpenConns(cfg.DBMaxOpenConnection)
	sqlMasterDB.SetMaxIdleConns(cfg.DBMaxIdleConnection)
	sqlMasterDB.SetConnMaxLifetime(cfg.DBConnMaxLife * time.Second)

	sqlReplicaDB.SetMaxOpenConns(cfg.DBMaxOpenConnection)
	sqlReplicaDB.SetMaxIdleConns(cfg.DBMaxIdleConnection)
	sqlReplicaDB.SetConnMaxLifetime(cfg.DBConnMaxLife * time.Second)

	// Test connections
	if err := sqlMasterDB.Ping(); err != nil {
		log.Fatalf("error pinging master database: %v", err)
	}

	if err := sqlReplicaDB.Ping(); err != nil {
		log.Fatalf("error pinging replica database: %v", err)
	}

	log.Printf("Database initialization successful")
	log.Printf("Master DB - Open connections: %d, Idle connections: %d",
		sqlMasterDB.Stats().OpenConnections, sqlMasterDB.Stats().Idle)
	log.Printf("Replica DB - Open connections: %d, Idle connections: %d",
		sqlReplicaDB.Stats().OpenConnections, sqlReplicaDB.Stats().Idle)

	masterDB.AutoMigrate(&model.MigrationRecord{})

	migrations.Migrate(context.Background(), masterDB)

	return masterDB, replicaDB
}

func GetReadDSN(c config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHostRead, c.DBPortRead, c.DBUserRead, c.DBPassword, c.DBName)
}

func GetWriteDSN(c config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHostWrite, c.DBPortWrite, c.DBUserWrite, c.DBPassword, c.DBName)
}
