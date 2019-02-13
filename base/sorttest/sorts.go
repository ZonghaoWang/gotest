package sorttest

import (
	"github.com/ZonghaoWang/gotest/base/convertindex"
	"math/rand"
)

const (
	Credit = iota
	Active
	FaceCount
)

type ProfileSlice struct {
	Profiles []*convertindex.Profile
	SortIndex int
}

func (p ProfileSlice) Len() int {
	return len(p.Profiles)
}

func (p ProfileSlice) Swap(i, j int) {
	p.Profiles[i], p.Profiles[j] = p.Profiles[j], p.Profiles[i]
}

func (p ProfileSlice) Less(i, j int) bool {
	switch p.SortIndex {
	case Credit:
		return p.Profiles[i].Credit > p.Profiles[j].Credit
	case Active:
		return p.Profiles[i].Active > p.Profiles[j].Active
	case FaceCount:
		return p.Profiles[i].FaceCount > p.Profiles[j].FaceCount
	default:
		return p.Profiles[i].Credit > p.Profiles[j].Credit
	}
}

func New(length int, randSeed int64, sortIndex int) ProfileSlice {
	var profileSlice ProfileSlice
	profileSlice.Profiles = make([]*convertindex.Profile, length)
	rand.Seed(randSeed)
	for index := 0; index < length; index++ {
		tmp := convertindex.Profile{
			Credit: rand.Int63n(10000),
			Active: rand.Int63n(10000),
			FaceCount: rand.Int63n(10000),
		}
		profileSlice.Profiles[index] = &tmp
	}
	profileSlice.SortIndex = sortIndex
	return profileSlice
}