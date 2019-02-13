package convertindex

import (
	"math/rand"
	"time"
	"sync"
	"sort"
	"fmt"
	"strconv"
)
const (
	Credit = iota
	Active
	FaceCount
)

const (
	MaxCreditBitShift = 16
	MaxActiveBitShift = 16
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
	FaceId		string
	CaptureTime time.Time
}

type Profile struct {
	ProfileID	string
	Credit 		int64
	Active 		int64
	FaceCount 	int64
	FaceList	[]*Face
	Match 		int64
}
//
//type IndexWeight struct {
//	CreditWeight 	[2 ^ MaxCreditBitShift]int64
//	ActiveWeight 	[2 ^ MaxActiveBitShift]int64
//	FaceCountWeight [2 ^ MaxFaceCountBitShift]int64
//	sync.Mutex
//}
//
//func (iw IndexWeight) Init()  {
//	for index, _ := range iw.CreditWeight {
//		iw.CreditWeight[index] = CreditCover + int64(index << CreditIndexBitShift)
//	}
//	for index, _ := range iw.ActiveWeight {
//		iw.ActiveWeight[index] = ActiveCover + int64(index << ActiveIndexBitShift)
//	}
//	for index, _ := range iw.FaceCountWeight {
//		iw.FaceCountWeight[index] = FaceCountCover + int64(index << FaceCountIndexBitShift)
//	}
//}
//func (iw IndexWeight) GetAndUpdate(index, value int) int64 {
//	switch index {
//	case Credit:
//		iw.CreditWeight[value]--
//		return iw.CreditWeight[value]
//	case Active:
//		iw.ActiveWeight[value]--
//		return iw.ActiveWeight[value]
//	case FaceCount:
//		iw.FaceCountWeight[value]--
//		return iw.FaceCountWeight[value]
//	default:
//		iw.CreditWeight[value]--
//		return iw.CreditWeight[value]
//
//	}
//}
//
//func (iw IndexWeight) Get(index, value int) int64 {
//	switch index {
//	case Credit:
//		return iw.CreditWeight[value]
//	case Active:
//		return iw.ActiveWeight[value]
//	case FaceCount:
//		return iw.FaceCountWeight[value]
//	default:
//		return iw.CreditWeight[value]
//	}
//}

type Memory struct {
	Profiles map[string]*Profile
	sync.RWMutex
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

type ProfileSliceCount struct {
	Profiles 	[]*Profile
	TotalNum	int
}

func (psc *ProfileSliceCount) AddProfile(profile *Profile) {
	psc.Profiles = append(psc.Profiles, profile)
	psc.TotalNum++
}

func (psc *ProfileSliceCount) RemoveProfile()  {
	psc.TotalNum--
}

func (psc *ProfileSliceCount) Update(ceil, total int) {
	psc.Profiles = psc.Profiles[0: ceil]
	psc.TotalNum = total
}
type InvertIndex struct {
	CreditIndex 	map[int64]*ProfileSliceCount
	ActiveIndex 	map[int64]*ProfileSliceCount
	FaceCountIndex 	map[int64]*ProfileSliceCount
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
		ii.ActiveIndex[ivfb.Credit].RemoveProfile()
	}
	if ivfa.ProfilePtr != nil {
		ii.ActiveIndex[ivfa.Credit].AddProfile(ivfa.ProfilePtr)
	}
}

func (ii *InvertIndex) UpdateFaceCount(ivfb, ivfa InvertFeature) {
	if ivfb.ProfilePtr != nil {
		ii.FaceCountIndex[ivfb.Credit].RemoveProfile()
	}
	if ivfa.ProfilePtr != nil {
		ii.FaceCountIndex[ivfa.Credit].AddProfile(ivfa.ProfilePtr)
	}
}


func (ii *InvertIndex) Update(ivfb, ivfa InvertFeature){
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
	memory 			*Memory
	invertIndex		*InvertIndex
}


func (m *Memory) Init() {
	m.Profiles = make(map[string]*Profile, leng)
	rand.Seed(47)
	for index := 0; index < leng; index++ {
		tmp := Profile{
			Credit: rand.Int63n(2 ^ MaxCreditBitShift),
			Active: rand.Int63n(2 ^ MaxActiveBitShift),
			FaceCount: rand.Int63n(2 ^ MaxFaceCountBitShift),
			ProfileID: strconv.Itoa(index),
			Match: rand.Int63n(10),
		}
		m.Profiles[strconv.Itoa(index)] = &tmp
	}
}

type ProfileSlice struct {
	Profiles []*Profile
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

const (
	MatchByInfoCost = 30
	MatchResultProp = 10
)
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
		Profiles: result,
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
		var result []*Profile
		profileListWapper := indexItem[key]
		profileList := profileListWapper.Profiles
		i, j := 0, 0
		for {
			if j == len(profileList) {
				break
			}
			if profileList[j] != nil && ((index == Credit && profileList[j].Credit == key) || (index == Active && profileList[j].Active == key) || (index == FaceCount && profileList[j].FaceCount == key)) && sp.memory.Profiles[profileList[j].ProfileID] != nil {
				profileList[i] = profileList[j]
				if MatchByInfo(profileList[i]) {
					result = append(result, profileList[i])
				}
				i++
			}
			j++
		}
		if i == profileListWapper.TotalNum {
			fmt.Printf("totalNum equals to valided profiles, number is %d\n", i)
			indexItem[key].Update(i, i)
			//indexItem[key] = &ProfileSliceCount{
			//	Profiles: profileList[0: i],
			//	TotalNum: i,
			//}
			results = append(results, result...)
		} else {
			// 应该有重复的profileId
			tmpSet := make(map[string]bool, i)
			var result2 []*Profile
			for inIndex, profile := range profileList[0: i] {
				if tmpSet[profile.ProfileID] {
					fmt.Printf("set have this profileid %s\n", profile.ProfileID)
					profileList[inIndex] = nil
				} else {
					tmpSet[profile.ProfileID] = true
					if MatchByInfo(profile){
						result2 = append(result2, profile)
					}
				}
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