package service

import (
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
)

type Clock interface {
	Now() time.Time
}

type BucketService struct {
	clock Clock
}

func NewBucketService(timeLocation time.Location) BucketService {
	return BucketService{
		clock: ClockWithTimeLocation{
			timeLocation: timeLocation,
		},
	}
}

func (b BucketService) CurrentBucket() model.Bucket {
	now := b.clock.Now()
	return b.generateBucket(now)
}

func (b BucketService) generateBucket(bucketTime time.Time) model.Bucket {
	return model.Bucket{
		Key: bucketTime.Format("2006-01-02"),
	}
}

// Clock implementation
type ClockWithTimeLocation struct {
	timeLocation time.Location
}

func (c ClockWithTimeLocation) Now() time.Time {
	return time.Now().In(&c.timeLocation)
}
