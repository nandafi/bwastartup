package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error) //tanpa parameter | feedback berupa slice Campaign(struct)
	FindByUserID(userID int) ([]Campaign, error)
}

//definisikan struct untuk akses Database
type repository struct {
	db *gorm.DB
}

//supaya bisa diakses dari luar package, buatkan instance/func
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign

	//pakai preload untuk menghubungkan relasi ke tabel campaign_images (struct CampaignImage)
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil

}

func (r *repository) FindByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}
