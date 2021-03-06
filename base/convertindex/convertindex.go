package convertindex

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	Credit = iota
	Active
	FaceCount
)


const (
	MatchByInfoCost = 3
	MatchResultProp = 10
)


const (
	MaxCreditBitShift    = 16
	MaxActiveBitShift    = 16
	MaxFaceCountBitShift = 16
	//CreditCover = 2 ^ (63 - MaxCreditBitShift)
	//ActiveCover = 2 ^ (63 - MaxActiveBitShift)
	//FaceCountCover = 2 ^ (63 - MaxFaceCountBitShift)
	//CreditIndexBitShift = 63 - MaxCreditBitShift
	//ActiveIndexBitShift = 63 - MaxActiveBitShift
	//FaceCountIndexBitShift = 63 - MaxFaceCountBitShift
	leng = 1000000
)

type Face struct {
	FaceId      string
	CaptureTime time.Time
}

type Profile struct {
	ProfileID string
	Credit    int64
	Active    int64
	FaceCount int64
	FaceList  []*Face
	Match     int64
}

type Memory struct {
	Profiles map[string]*Profile
	sync.RWMutex
}

type InvertFeature struct {
	Credit     int64
	Active     int64
	FaceCount  int64
	ProfilePtr *Profile
}

func CacheFeature(profile *Profile) (result InvertFeature) {
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

type ProfileSliceCount struct {
	Profiles []string
	TotalNum int
}

func (psc *ProfileSliceCount) AddProfile(profile *Profile) {
	psc.Profiles = append(psc.Profiles, profile.ProfileID)
	psc.TotalNum++
}

func (psc *ProfileSliceCount) RemoveProfile() {
	psc.TotalNum--
}

func (psc *ProfileSliceCount) Update(ceil, total int) {
	psc.Profiles = psc.Profiles[0:ceil]
	psc.TotalNum = total
}

type InvertIndex struct {
	CreditIndex    map[int64]*ProfileSliceCount
	ActiveIndex    map[int64]*ProfileSliceCount
	FaceCountIndex map[int64]*ProfileSliceCount
	sync.Mutex
}

func (ii *InvertIndex) UpdateCredit(ivfb, ivfa InvertFeature) {
	if ivfb.ProfilePtr != nil {
		ii.CreditIndex[ivfb.Credit].RemoveProfile()
	}
	if ivfa.ProfilePtr != nil {
		ii.CreditIndex[ivfa.Credit].AddProfile(ivfa.ProfilePtr)
	}
}

func (ii *InvertIndex) UpdateActive(ivfb, ivfa InvertFeature) {
	if ivfb.ProfilePtr != nil {
		ii.ActiveIndex[ivfb.Active].RemoveProfile()
	}
	if ivfa.ProfilePtr != nil {
		ii.ActiveIndex[ivfa.Active].AddProfile(ivfa.ProfilePtr)
	}
}

func (ii *InvertIndex) UpdateFaceCount(ivfb, ivfa InvertFeature) {
	if ivfb.ProfilePtr != nil {
		ii.FaceCountIndex[ivfb.FaceCount].RemoveProfile()
	}
	if ivfa.ProfilePtr != nil {
		ii.FaceCountIndex[ivfa.FaceCount].AddProfile(ivfa.ProfilePtr)
	}
}

func (ii *InvertIndex) Update(ivfb, ivfa InvertFeature) {
	if !(ivfb.ProfilePtr != nil && ivfa.ProfilePtr != nil && ivfb.Credit == ivfa.Credit) {
		ii.UpdateCredit(ivfb, ivfa)
	}
	if !(ivfb.ProfilePtr != nil && ivfa.ProfilePtr != nil && ivfb.Active == ivfa.Active) {
		ii.UpdateActive(ivfb, ivfa)
	}
	if !(ivfb.ProfilePtr != nil && ivfa.ProfilePtr != nil && ivfb.FaceCount == ivfa.FaceCount) {
		ii.UpdateFaceCount(ivfb, ivfa)
	}
}

type SuperProfile struct {
	memory      *Memory
	invertIndex *InvertIndex
}

func (m *Memory) Init() {
	m.Profiles = make(map[string]*Profile, leng)
	rand.Seed(47)
	for index := 0; index < leng; index++ {
		tmp := Profile{
			Credit:    rand.Int63n(2 << MaxCreditBitShift - 1),
			Active:    rand.Int63n(2 << MaxActiveBitShift - 1),
			FaceCount: rand.Int63n(2 << MaxFaceCountBitShift - 1),
			ProfileID: strconv.Itoa(index),
			Match:     rand.Int63n(MatchResultProp),
		}
		m.Profiles[strconv.Itoa(index)] = &tmp
	}
}

type ProfileSlice struct {
	Profiles  []*Profile
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
func (ii *InvertIndex) Init(memory *Memory) {
	//memory.Lock()
	ii.CreditIndex = make(map[int64]*ProfileSliceCount)
	ii.ActiveIndex = make(map[int64]*ProfileSliceCount)
	ii.FaceCountIndex = make(map[int64]*ProfileSliceCount)
	for _, profile := range memory.Profiles {
		credit := profile.Credit
		active := profile.Active
		faceCount := profile.FaceCount
		if ii.CreditIndex[credit] == nil {
			ii.CreditIndex[credit] = &ProfileSliceCount{}
		}
		if ii.ActiveIndex[active] == nil {
			ii.ActiveIndex[active] = &ProfileSliceCount{}
		}
		if ii.FaceCountIndex[faceCount] == nil {
			ii.FaceCountIndex[faceCount] = &ProfileSliceCount{}
		}
		ii.CreditIndex[credit].AddProfile(profile)
		ii.ActiveIndex[active].AddProfile(profile)
		ii.FaceCountIndex[faceCount].AddProfile(profile)
	}
	//fmt.Println(ii.ActiveIndex[1])
	//memory.Unlock()
}

func (sp *SuperProfile) Init() {
	fmt.Println("init memory...")
	sp.memory.Init()
	fmt.Println("init invertIndex...")
	sp.invertIndex.Init(sp.memory)
}



func MatchByInfo(profile *Profile) bool {
	time.Sleep(time.Duration(rand.Int63n(MatchByInfoCost)) * time.Nanosecond)
	return profile.Match < 10
}

func (sp *SuperProfile) SearchByInfo(index int) []*Profile {
	var result []*Profile
	sp.memory.RLock()
	for _, profile := range sp.memory.Profiles {
		if MatchByInfo(profile) {
			result = append(result, profile)
		}
	}
	sort.Sort(ProfileSlice{
		Profiles:  result,
		SortIndex: index,
	})
	sp.memory.RUnlock()
	return result
}

func (sp *SuperProfile) SearchByInfoIndex(index int) []*Profile {
	sp.invertIndex.Lock()
	var keys []int64
	var indexItem map[int64]*ProfileSliceCount
	switch index {
	case Credit:
		indexItem = sp.invertIndex.CreditIndex
	case Active:
		indexItem = sp.invertIndex.ActiveIndex
	case FaceCount:
		indexItem = sp.invertIndex.FaceCountIndex
	default:
		indexItem = sp.invertIndex.CreditIndex
	}

	for key := range indexItem {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	var results []*Profile
	for _, key := range keys {
		result := make([]*Profile, 0, len(indexItem[key].Profiles) / 2)
		profileListWapper := indexItem[key]
		profileList := profileListWapper.Profiles
		i, j := 0, 0
		sp.memory.RLock()
		for {
			if j == len(profileList) {
				break
			}
			if profileList[j] == "" {
				j++
				continue
			}
			tmpProfile := sp.memory.Profiles[profileList[j]]
			if tmpProfile != nil && ((index == Credit && tmpProfile.Credit == key) || (index == Active && tmpProfile.Active == key) || (index == FaceCount && tmpProfile.FaceCount == key)){
				profileList[i] = profileList[j]
				if MatchByInfo(tmpProfile) {
					result = append(result, tmpProfile)
				}
				i++
			}
			j++
		}
		sp.memory.RUnlock()
		if i == profileListWapper.TotalNum {
			//fmt.Printf("totalNum equals to valided profiles, number is %d\n", i)
			indexItem[key].Update(i, i)
			results = append(results, result...)
		} else {
			// 应该有重复的profileId
			//fmt.Println("real not equals to fact")
			tmpSet := make(map[string]bool, i)
			var result2 []*Profile
			for inIndex, profileId := range profileList[0:i] {
				sp.memory.RLock()
				if tmpSet[profileId] {
					//fmt.Printf("set have this profileid %s\n", profileId)
					profileList[inIndex] = ""
				} else {
					tmpSet[profileId] = true
					tmpProfile := sp.memory.Profiles[profileId]
					if MatchByInfo(tmpProfile) {
						result2 = append(result2, tmpProfile)
					}
				}
				sp.memory.RUnlock()
			}
			indexItem[key].Update(i, len(tmpSet))
			//indexItem[key] = &ProfileSliceCount{
			//	Profiles: profileList[0: i],
			//	TotalNum: len(tmpSet),
			//}
			results = append(results, result2...)
		}
	}
	sp.invertIndex.Unlock()
	return results
}
