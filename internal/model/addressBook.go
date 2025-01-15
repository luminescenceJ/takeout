package model

type AddressBook struct {
	CityCode     string `json:"cityCode,omitempty"`
	CityName     string `json:"cityName,omitempty"`
	Consignee    string `json:"consignee,omitempty"`
	Detail       string `json:"detail"` // 详细地址
	DistrictCode string `json:"districtCode,omitempty"`
	DistrictName string `json:"districtName,omitempty"`
	ID           int64  `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	IsDefault    int64  `json:"isDefault,omitempty"`
	Label        string `json:"label,omitempty"`
	Phone        string `json:"phone"` // 手机号
	ProvinceCode string `json:"provinceCode,omitempty"`
	ProvinceName string `json:"provinceName,omitempty"`
	Sex          string `json:"sex"`
	UserID       int64  `json:"userId,omitempty"`
}
