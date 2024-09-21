package service

import (
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
)

type BucketService struct {
	timeLocation *time.Location
}

func NewBucketService(timeLocation *time.Location) BucketService {
	return BucketService{
		timeLocation: timeLocation,
	}
}

func (b BucketService) CurrentBucket() model.Bucket {
	now := time.Now().In(b.timeLocation)
	return b.generateBucket(now)
}

func (b BucketService) generateBucket(bucketTime time.Time) model.Bucket {
	return model.Bucket{
		Key: bucketTime.Format("2006-01-02"),
	}
}
