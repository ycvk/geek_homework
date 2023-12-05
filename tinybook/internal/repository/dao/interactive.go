package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Interactive struct {
	Id    int64  `gorm:"column:id;primaryKey;autoIncrement;not null"`
	BizId int64  `gorm:"column:biz_id;not null;uniqueIndex:idx_biz_type_id"`
	Biz   string `gorm:"column:biz;not null;type:varchar(32);uniqueIndex:idx_biz_type_id"`

	ReadCount    int64 `gorm:"column:read_count"`
	LikeCount    int64 `gorm:"column:like_count"`
	CollectCount int64 `gorm:"column:collect_count"`
	Utime        int64 `gorm:"column:utime;not null"`
	Ctime        int64 `gorm:"column:ctime;not null"`
}

type LikeRecord struct {
	Id     int64  `gorm:"column:id;primaryKey;autoIncrement;not null"`
	Uid    int64  `gorm:"column:uid;not null;uniqueIndex:idx_biz_like_id"`
	BizId  int64  `gorm:"column:biz_id;not null;uniqueIndex:idx_biz_like_id"`
	Biz    string `gorm:"column:biz;not null;type:varchar(32);uniqueIndex:idx_biz_like_id"`
	Status uint8  `gorm:"column:status;not null;type:tinyint(1)"`
	Utime  int64  `gorm:"column:utime;not null"`
	Ctime  int64  `gorm:"column:ctime;not null"`
}

type InteractiveDAO interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
	InsertLikeRecord(ctx context.Context, biz string, id int64, uid int64) error
	DeleteLikeRecord(ctx context.Context, biz string, id int64, uid int64) error
}

type GormInteractiveDAO struct {
	db *gorm.DB
}

func NewGormInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GormInteractiveDAO{db: db}
}

func (g *GormInteractiveDAO) InsertLikeRecord(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().Unix()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先插入点赞记录
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"utime":  now,
				"status": 1,
			}),
		}).Create(&LikeRecord{
			BizId:  id,
			Biz:    biz,
			Uid:    uid,
			Utime:  now,
			Ctime:  now,
			Status: 1,
		}).Error
		// 如果有错误就回滚
		if err != nil {
			return err
		}
		// 然后更新点赞数
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_count": gorm.Expr("like_count + ?", 1),
				"utime":      now,
			}),
		}).Create(&Interactive{
			BizId:     id,
			Biz:       biz,
			LikeCount: 1,
			Utime:     now,
			Ctime:     now,
		}).Error
	})
}

func (g *GormInteractiveDAO) DeleteLikeRecord(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().Unix()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&LikeRecord{}).
			Where("uid = ? and biz_id = ? and biz = ?", uid, id, biz).
			Updates(map[string]interface{}{
				"status": 0,
				"utime":  now,
			}).Error
		if err != nil {
			return err
		}
		return tx.Model(&Interactive{}).
			Where("biz_id = ? and biz = ?", id, biz).
			Updates(map[string]interface{}{
				"like_count": gorm.Expr("like_count - ?", 1),
				"utime":      now,
			}).Error
	})
}

func (g *GormInteractiveDAO) IncreaseReadCount(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().Unix()
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_count": gorm.Expr("read_count + ?", 1),
			"utime":      now,
		}),
	}).Create(&Interactive{
		BizId:     bizId,
		Biz:       biz,
		ReadCount: 1,
		Ctime:     now,
		Utime:     now,
	}).Error
}
