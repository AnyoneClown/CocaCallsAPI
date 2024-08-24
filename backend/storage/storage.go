package storage

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/AnyoneClown/CocaCallsAPI/types"
    "github.com/google/uuid"
    "github.com/joho/godotenv"
    "github.com/go-gormigrate/gormigrate/v2"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type CockroachDB struct {
    db *gorm.DB
}

func NewCockroachDB() *CockroachDB {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("error loading .env file")
    }

    dsn := os.Getenv("COCKROACH_DB_URL")
    if dsn == "" {
        log.Fatalf("COCKROACH_DB_URL environment variable not set")
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }

    if err := runMigrations(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    return &CockroachDB{db: db}
}

func runMigrations(db *gorm.DB) error {
    log.Println("Starting migrations...")

    m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
        {
            ID: "202408241000",
            Migrate: func(tx *gorm.DB) error {
                log.Println("Running migration 202408241000...")
                return tx.AutoMigrate(&types.User{})
            },
            Rollback: func(tx *gorm.DB) error {
                log.Println("Rolling back migration 202408241000...")
                return tx.Migrator().DropTable("users")
            },
        },
    })

    if err := m.Migrate(); err != nil {
        return fmt.Errorf("could not migrate: %v", err)
    }

    log.Println("Migration did run successfully")
    return nil
}

func (c *CockroachDB) CreateUser(email, password, googleID, picture, provider string, verifiedEmail bool) (types.User, error) {
    if provider == "" {
        if err := types.ValidateUser(email, password); err != nil {
            return types.User{}, err
        }
    } else {
        if googleID == "" || email == "" {
            return types.User{}, fmt.Errorf("invalid OAuth user data")
        }
    }

    var existingUser types.User
    if err := c.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
        return types.User{}, fmt.Errorf("email already in use")
    } else if err != gorm.ErrRecordNotFound {
        return types.User{}, err
    }

    if provider != "" {
        if err := c.db.Where("google_id = ?", googleID).First(&existingUser).Error; err == nil {
            return types.User{}, fmt.Errorf("Google ID already in use")
        } else if err != gorm.ErrRecordNotFound {
            return types.User{}, err
        }
    }

    user := types.User{
        ID:            uuid.New(),
        Email:         email,
        Password:      password,
        GoogleID:      googleID,
        Picture:       picture,
        Provider:      provider,
        VerifiedEmail: verifiedEmail,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    result := c.db.Create(&user)
    if result.Error != nil {
        return types.User{}, result.Error
    }

    return user, nil
}

func (c *CockroachDB) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	result := c.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return types.User{}, fmt.Errorf("user not found")
		}
		return types.User{}, result.Error
	}

	return user, nil
}
