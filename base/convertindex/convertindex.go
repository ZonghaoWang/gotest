package convertindex

import (
	"math/rand"
	"time"
)

var leng = 10000000

type Face struct {
	FaceId		string
	CaptureTime time.Time
}

type Profile struct {
	Credit 		int64
	Active 		int64
	FaceCount 	int64
	FaceList	[]*Face
}

type Memory struct {
	Profiles []*Profile

}

type InvertFeature struct {
	Credit		int64
	Active 		int64
	FaceCount 	int64
	ProfilePtr 	*Profile
}

func (ivf *InvertIndex) CacheFeature(profile *Profile) (result InvertFeature) {
	if profile != nil {
		result = InvertFeature{
			Credit:     profile.Credit,
			Active:     profile.Active,
			FaceCount:  profile.FaceCount,
			ProfilePtr: profile,
		}
	}
	return
}

type InvertIndex struct {
	CreditIndex 	map[int64][]*Profile
	ActiveIndex 	map[int64][]*Profile
	FaceCountIndex 	map[int64][]*Profile
}

func (ii *InvertIndex) UpdateCredit(profile *Profile) {
	if value, exist := ii.CreditIndex[profile.Active]; exist {
		ii.ActiveIndex[profile.Active] = append(value, profile)
	} else {
		ii.ActiveIndex[profile.Active] = []*Profile{profile}
	}
}

func (ii *InvertIndex) UpdateActive(profile *Profile) {
	if value, exist := ii.ActiveIndex[profile.Credit]; exist {
		ii.CreditIndex[profile.Credit] = append(value, profile)
	} else {
		ii.CreditIndex[profile.Credit] = []*Profile{profile}
	}
}

func (ii *InvertIndex) UpdateFaceCount(profile *Profile) {
	if value, exist := ii.FaceCountIndex[profile.FaceCount]; exist {
		ii.FaceCountIndex[profile.FaceCount] = append(value, profile)
	} else {
		ii.FaceCountIndex[profile.FaceCount] = []*Profile{profile}
	}
}


func (ii *InvertIndex) Update(ivf InvertFeature, profile *Profile){
	var creditIndex, activeIndex, faceCountIndex int
	if ivf.ProfilePtr == nil || ivf.Credit == profile.Credit {
		creditIndex
	}
}


func (m *Memory) Init() {
	m.Profiles = make([]*Profile, leng)
	rand.Seed(47)
	for index := 0; index < leng; index++ {
		tmp := Profile{
			Credit: rand.Int63n(10000),
			Active: rand.Int63n(10000),
			FaceCount: rand.Int63n(10000),
		}
		m.Profiles[index] = &tmp
	}
}