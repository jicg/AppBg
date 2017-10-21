package db

import (
	log "github.com/sirupsen/logrus"
)

type Offorder struct {
	Model
	Orderno  string  `json:"orderno" gorm:"unique_index"`
	Billdate string  `json:"billdate"`
	Optname  string  `json:"optname"`
	Remark   string  `json:"remark"`
	Amt      float32 `json:"amt"`
	Status   uint8   `json:"status"`
	Product  string  `json:"product"`
	Sttype   uint8   `json:"sttype"`
}

func QueryOfforder(order string, datebeg string, dateend string, optname string, remark string, product string, sttype uint8, status uint8, start int, size int) ([]*Offorder, int) {
	orders := make([]*Offorder, size)
	count := 0
	tempdb := db;
	if len(order) != 0 {
		tempdb = tempdb.Where("orderno = ?", order)
	}
	if len(optname) != 0 {
		tempdb = tempdb.Where("optname = ?", optname)
	}
	if len(remark) != 0 {
		tempdb = tempdb.Where("remark like ?", "%"+remark+"%")
	}
	if len(product) != 0 {
		tempdb = tempdb.Where("product like ?", "%"+product+"%")
	}
	if status != 0 {
		tempdb = tempdb.Where("status = ?", status)
	}
	if sttype != 0 {
		tempdb = tempdb.Where("sttype = ?", sttype)
	}
	if len(datebeg) != 0 {
		tempdb = tempdb.Where("billdate >= ?", datebeg)
	}
	if len(dateend) != 0 {
		tempdb = tempdb.Where("billdate <= ?", dateend)
	}
	if err := tempdb.Order("id desc").Offset(start).Limit(size).Find(&orders).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		log.Error(err)
		return orders, 0
	}
	return orders, count
}


func QueryOfforderWithOutPage(order string, datebeg string, dateend string, optname string, remark string, product string, sttype uint8, status uint8) ([]*Offorder) {
	orders := make([]*Offorder, 0)
	tempdb := db;
	if len(order) != 0 {
		tempdb = tempdb.Where("orderno = ?", order)
	}
	if len(optname) != 0 {
		tempdb = tempdb.Where("optname = ?", optname)
	}
	if len(remark) != 0 {
		tempdb = tempdb.Where("remark like ?", "%"+remark+"%")
	}
	if len(product) != 0 {
		tempdb = tempdb.Where("product like ?", "%"+product+"%")
	}
	if status != 0 {
		tempdb = tempdb.Where("status = ?", status)
	}
	if sttype != 0 {
		tempdb = tempdb.Where("sttype = ?", sttype)
	}
	if len(datebeg) != 0 {
		tempdb = tempdb.Where("billdate >= ?", datebeg)
	}
	if len(dateend) != 0 {
		tempdb = tempdb.Where("billdate <= ?", dateend)
	}
	if err := tempdb.Order("id desc").Find(&orders).Error; err != nil {
		log.Error(err)
		return orders
	}
	return orders
}


func GetOfforderById(id int) (*Offorder, error) {
	offorder := &Offorder{}
	if err := db.Model(offorder).Where("id = ?", id).Limit(1).Scan(offorder).Error; err != nil {
		return nil, err
	}
	return offorder, nil
}

func GetOfforderByOrderno(orderno string) (*Offorder, error) {
	offorder := &Offorder{}
	if err := db.Model(offorder).Where("orderno = ?", orderno).Limit(1).Scan(offorder).Error; err != nil {
		return nil, err
	}
	return offorder, nil
}

func AddOfforder(offorder *Offorder) error {
	return db.Create(offorder).Error
}

//func UpdateOfforder(offorder *Offorder) error {
//	return db.Updates(offorder).Error
//}

func UpdateOfforder(offorder *Offorder) (error) {
	return db.Model(&Offorder{}).Update(offorder).Error
}

func DeleteOfforder(id int64) (error) {
	return db.Where("id = ?", id).Delete(&Offorder{}).Error
}
