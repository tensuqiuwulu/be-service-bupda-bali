package repository

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"gorm.io/gorm"
)

type DesaRepositoryInterface interface {
	FindDesaByIdKelu(db *gorm.DB, idKelu int) ([]entity.Desa, error)
}

type DesaRepositoryImplementation struct {
	DB *config.Database
}

func NewDesaRepository(
	db *config.Database,
) DesaRepositoryInterface {
	return &DesaRepositoryImplementation{
		DB: db,
	}
}

func (service *DesaRepositoryImplementation) FindDesaByIdKelu(db *gorm.DB, idKelu int) ([]entity.Desa, error) {
	desas := []entity.Desa{}
	results := db.
		Where("id_kelurahan = ?", idKelu).
		Find(&desas)
	return desas, results.Error
}
