package model

type AddressBook struct {
	ID           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64  `gorm:"not null" json:"userId"`
	Consignee    string `gorm:"size:50" json:"consignee"`
	Sex          string `gorm:"size:2" json:"sex"`
	Phone        string `gorm:"size:11;not null" json:"phone"`
	ProvinceCode string `gorm:"size:12" json:"provinceCode"`
	ProvinceName string `gorm:"size:32" json:"provinceName"`
	CityCode     string `gorm:"size:12" json:"cityCode"`
	CityName     string `gorm:"size:32" json:"cityName"`
	DistrictCode string `gorm:"size:12" json:"districtCode"`
	DistrictName string `gorm:"size:32" json:"districtName"`
	Detail       string `gorm:"size:200" json:"detail"`
	Label        string `gorm:"size:100" json:"label"`
	IsDefault    int8   `gorm:"default:0;not null" json:"isDefault"`
}
