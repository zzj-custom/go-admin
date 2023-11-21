package db

import (
	"github.com/zzj-custom/pkg/pMysql"
	"gorm.io/gorm"
	"log/slog"
	"sync"
	"time"
)

type ImagesModel struct {
	Id            int       `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Name          string    `gorm:"column:name;NOT NULL;comment:'图片名称'"`
	Copyright     string    `gorm:"column:copyright;default:;NOT NULL;comment:'版权'"`
	CopyrightLink string    `gorm:"column:copyright_link;default:;NOT NULL;comment:'版权链接'"`
	Url           string    `gorm:"column:url;NOT NULL;comment:'图片地址'"`
	Start         string    `gorm:"column:start;NOT NULL;comment:'开始时间'"`
	End           string    `gorm:"column:end;NOT NULL;comment:'结束时间'"`
	Location      string    `gorm:"column:location;default:zh-cn;NOT NULL;comment:'位置，中国:zh-CN'"`
	ClickCount    int       `gorm:"column:click_count;default:0;NOT NULL;comment:'点击次数'"`
	DownloadCount int       `gorm:"column:download_count;default:0;NOT NULL;comment:'下载次数'"`
	Hash          string    `gorm:"column:hash;NOT NULL;comment:'唯一hash值'"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;comment:'更新时间'"`
}

func (im *ImagesModel) TableName() string {
	return "bing_images"
}

type ImagesModelRepository struct {
	db *gorm.DB
}

var (
	imagesModelRepo *ImagesModelRepository
	imagesModelOnce sync.Once
)

func NewImagesModelRepository(db *gorm.DB) *ImagesModelRepository {
	imagesModelOnce.Do(func() {
		conn, err := pMysql.GetDb("go")
		if err != nil {
			slog.With("db", "go", "err", err).Error("Failed to get database")
			return
		}
		imagesModelRepo.db = conn
	})
	return imagesModelRepo
}
