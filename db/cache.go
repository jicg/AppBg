package db

type Cache struct {
	Model
	Name  string `json:"name" gorm:"unique_index"`
	Value string `json:"value"`
}

func GetCache(name string) (string, error) {
	var cache Cache
	if err := db.Model(&Cache{}).Where("name = ?", name).Limit(1).Scan(&cache).Error; err != nil {
		return "", err
	}
	return cache.Value, nil
}

func SaveCache(name string, value string) error {
	cnt := 0
	cache := &Cache{Name: name, Value: value}
	db.Model(&Cache{}).Where("name = ?", name).Count(&cnt);
	if cnt != 0 {
		db.Model(&Cache{}).Where("name = ?", name).Limit(1).Scan(cache)
		cache.Value = value
	}
	return db.Save(cache).Error
}
