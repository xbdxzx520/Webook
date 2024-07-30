package dao

import "context"

type AccountDAO interface {
	AddActivities(ctx context.Context, activities ...AccountActivity) error
}

// Account 账号本体
// 包含了当前状态
type Account struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	// 我账号是哪个用户的账号
	Uid int64

	// 唯一标识一个账号
	Account int64 `gorm:"uniqueIndex:account_type"`
	Type    uint8 `gorm:"uniqueIndex:account_type"`

	Balance  int64
	Currency string

	Utime int64
	Ctime int64
}

// AccountAudit, AccountBank...

type AccountActivity struct {
	Id  int64 `gorm:"primaryKey,autoIncrement"`
	Uid int64

	// 在 biz, biz_id, account 和 account_id 上创建一个联合唯一索引
	// 这样可以确保记账的时候不会重复记账
	Biz   string `gorm:"uniqueIndex:biz_type_id"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id"`

	Account     int64 `gorm:"index:account_type;uniqueIndex:biz_type_id"`
	AccountType uint8 `gorm:"index:account_type;uniqueIndex:biz_type_id"`

	// TYPE 入账还是出账
	Amount   int64
	Currency string

	Utime int64
	Ctime int64
}

func (AccountActivity) TableName() string {
	return "account_activities"
}
