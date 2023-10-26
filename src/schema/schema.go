package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Activity string

const (
	Pornography Activity = "Pornography"
	Drugs       Activity = "Drugs"
	Weaponry    Activity = "Weaponry"
)

type URL struct {
	ID            primitive.ObjectID `bson:"_id"`
	URL           string             `bson:"url"`
	IsSuspicious  bool               `bson:"isSuspicious"`
	Data          string             `bson:"data"`
	SusScore      int                `bson:"susScore"` //Suspicion score
	LastCrawled   time.Time          `bson:"lastCrawled"`
	IsOnline      bool               `bson:"isOnline"`
	CrawlCount    int                `bson:"crawlCount"`
	Type          []Activity         `bson:"type,omitempty"` //Types of illegal activities
	Links         []string           `bson:"links"`          //Links to other sites
	ClearnetSites []string           `bson:"clearnetSites"`  //Clearnet sites hosted by the same ip address
}

type UnscrapedURL struct {
	ID  primitive.ObjectID `bson:"_id"`
	URL string             `bson:"url"`
}
