package repo

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ProductRepo interface {
	Create(product *model.AutoPart) error
	GetAll(ctx context.Context) ([]model.AutoPart, string, error)
	GetByID(id uint) (*model.AutoPart, error)
	Update(product model.AutoPart, fieldsToUpdate map[string]interface{}) error
	Delete(id uint) error
	GetStock(id uint) (uint, error)
	UpdateStock(id uint, quantity int) error
	Exists(id uint) (bool, error)
	Search(query string) ([]model.AutoPart, error)
	FindByModelAndYear(modelName string, year int) ([]model.AutoPart, error)
	InvalidateCache() error
}

type productRepo struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewProductRepo(db *gorm.DB, redisClient *redis.Client) ProductRepo {
	return &productRepo{db: db, redisClient: redisClient}
}

func (ar *autoPartRepo) Create(autoPart *model.AutoPart) error {
	return ar.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(autoPart).Error; err != nil {
			return err
		}

		cacheKey := "all_auto_parts"
		var autoParts []model.AutoPart
		if err := tx.Preload("Categories").Preload("Brand").Preload("AutoPartInfo").Find(&autoParts).Error; err != nil {
			return err
		}

		autoPartsJSON, err := json.Marshal(&autoParts)
		if err != nil {
			return err
		}

		err = ar.redisClient.Set(context.Background(), cacheKey, autoPartsJSON, time.Hour).Err()
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *autoPartRepo) GetAll(ctx context.Context) ([]model.AutoPart, string, error) {
	cacheKey := "all_auto_parts"

	cachedAutoParts, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		var autoParts []model.AutoPart
		if err := r.db.Preload("Categories").Preload("Brand").Preload("AutoPartInfo").Find(&autoParts).Error; err != nil {
			return nil, "", err
		}

		autoPartsJSON, err := json.Marshal(&autoParts)
		if err != nil {
			return nil, "", err
		}
		err = r.redisClient.Set(ctx, cacheKey, autoPartsJSON, time.Hour).Err()
		if err != nil {
			return nil, "", err
		}

		log.Println("Data fetched from DB and cached")
		return autoParts, "DB", nil
	} else if err != nil {
		return nil, "", err
	}

	var autoParts []model.AutoPart
	err = json.Unmarshal([]byte(cachedAutoParts), &autoParts)
	if err != nil {
		return nil, "", err
	}

	log.Println("Data fetched from cache")
	return autoParts, "Cache", nil
}

func (r *autoPartRepo) GetByID(id uint) (*model.AutoPart, error) {
	var part model.AutoPart
	if err := r.db.Preload("Categories").Preload("Brand").Preload("AutoPartInfo").Where("id = ?", id).Find(&part).Error; err != nil {
		return nil, err
	}

	return &part, nil
}

func (ar *autoPartRepo) Update(product model.AutoPart, fieldsToUpdate map[string]interface{}) error {
	var existingProduct model.AutoPart
	if err := ar.db.Preload("Categories").First(&existingProduct, product.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("record not found")
		}
		return err
	}

	if categoryIDs, ok := fieldsToUpdate["category_id"].([]interface{}); ok {
		if err := ar.db.Model(&existingProduct).Association("Categories").Clear(); err != nil {
			return err
		}

		for _, id := range categoryIDs {
			catID, ok := id.(float64)
			if !ok {
				return errors.New("invalid category ID format")
			}

			var category model.Category
			if err := ar.db.First(&category, uint(catID)).Error; err != nil {
				return err
			}
			if err := ar.db.Model(&existingProduct).Association("Categories").Append(&category); err != nil {
				return err
			}
		}
		delete(fieldsToUpdate, "category_id")
	}

	if brandID, ok := fieldsToUpdate["brand_id"].(float64); ok {
		existingProduct.BrandID = uint(brandID)
		delete(fieldsToUpdate, "brand_id")
	}

	if forYears, ok := fieldsToUpdate["for_years"]; ok {
		forYearsJSON, err := json.Marshal(forYears)
		if err != nil {
			return errors.New("invalid JSON format for for_years")
		}
		fieldsToUpdate["for_years"] = string(forYearsJSON)
	}

	if len(fieldsToUpdate) > 0 {
		if err := ar.db.Model(&existingProduct).Updates(fieldsToUpdate).Error; err != nil {
			return err
		}
	}

	cacheKey := "all_auto_parts"
	if err := ar.redisClient.Del(context.Background(), cacheKey).Err(); err != nil {
		return err
	}

	return nil
}

func (ar *autoPartRepo) Delete(id uint) error {
	return ar.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("auto_part_id = ?", id).Delete(model.AutoPartCategory{}).Error; err != nil {
			return err
		}

		if err := tx.Where("auto_part_id = ?", id).Delete(model.ModelAutoAutoPart{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.AutoPart{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (ar *autoPartRepo) UpdateStock(id uint, quantity int) error {
	return ar.db.Model(&model.AutoPart{}).Where("id = ?", id).Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (ar *autoPartRepo) GetStock(id uint) (uint, error) {
	var autoPart model.AutoPart
	if err := ar.db.Select("stock").Where("id = ?", id).First(&autoPart).Error; err != nil {
		return 0, err
	}
	return autoPart.Stock, nil
}

func (r *autoPartRepo) Exists(id uint) (bool, error) {
	var count int64
	if err := r.db.Model(&model.AutoPart{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *autoPartRepo) Search(query string) ([]model.AutoPart, error) {
	var products []model.AutoPart

	terms := strings.Fields(query)
	likeClauses := make([]string, len(terms))
	likeValues := make([]interface{}, len(terms)*7)

	for i, term := range terms {
		likeClauses[i] = `(auto_parts.name ILIKE ? OR categories.name ILIKE ? OR brands.name ILIKE ? OR auto_parts.model_name ILIKE ? OR model_autos.name ILIKE ? OR model_autos.release_year @> ?::jsonb OR auto_parts.for_years @> ?::jsonb)`
		for j := 0; j < 5; j++ {
			likeValues[i*7+j] = "%" + term + "%"
		}
		likeValues[i*7+5] = `"` + term + `"`
		likeValues[i*7+6] = `"` + term + `"`
	}

	whereClause := strings.Join(likeClauses, " AND ")

	if err := r.db.Table("auto_parts").
		Select("auto_parts.*, array_agg(categories.name) AS category_names, brands.name AS brand_name").
		Joins("LEFT JOIN auto_part_categories ON auto_parts.id = auto_part_categories.auto_part_id").
		Joins("LEFT JOIN categories ON auto_part_categories.category_id = categories.id").
		Joins("LEFT JOIN brands ON auto_parts.brand_id = brands.id").
		Joins("LEFT JOIN model_auto_auto_parts ON auto_parts.id = model_auto_auto_parts.auto_part_id").
		Joins("LEFT JOIN model_autos ON model_auto_auto_parts.model_auto_id = model_autos.id").
		Where(whereClause, likeValues...).
		Group("auto_parts.id, brands.name").
		Find(&products).Error; err != nil {
		log.Printf("Error searching auto parts: %v", err)
		return nil, err
	}
	return products, nil

}

func (ar *autoPartRepo) FindByModelAndYear(modelName string, year int) ([]model.AutoPart, error) {
	var autoParts []model.AutoPart

	err := ar.db.
		Table("auto_parts").
		Select("auto_parts.*, model_autos.name AS model_name").
		Joins("JOIN model_auto_auto_parts ON auto_parts.id = model_auto_auto_parts.auto_part_id").
		Joins("JOIN model_autos ON model_auto_auto_parts.model_auto_id = model_autos.id").
		Where("auto_parts.for_years @> ?", []int{year}).
		Where("model_autos.name = ?", modelName).
		Find(&autoParts).Error

	if err != nil {
		log.Printf("Error finding auto parts by model and year: %v", err)
		return nil, err
	}

	return autoParts, nil
}

func (ar *autoPartRepo) InvalidateCache() error {
	cacheKey := "all_auto_parts"
	if err := ar.redisClient.Del(context.Background(), cacheKey).Err(); err != nil {
		return err
	}

	return nil
}
