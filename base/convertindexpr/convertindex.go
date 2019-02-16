package convertindexpr

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	Credit = iota
	Active
	FaceCount
)


const (
	MatchByInfoCost = 100
	MatchResultProp = 100
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
	Profiles 	[]*Profile
	TotalNum 	int
	RemoveFlag 	bool
}

func (psc *ProfileSliceCount) Reset(profiles []*Profile) {
	psc.Profiles = profiles
	psc.RemoveFlag = false
}

func (psc *ProfileSliceCount) AddProfile(profile *Profile) {
	psc.Profiles = append(psc.Profiles, profile)
	psc.TotalNum++
}

func (psc *ProfileSliceCount) RemoveProfile() {
	psc.TotalNum--
	psc.RemoveFlag = true
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
	if ii.CreditIndex[ivfa.Credit] == nil {
		ii.CreditIndex[ivfa.Credit] = &ProfileSliceCount{}
	}
	if ivfa.ProfilePtr != nil {
		ii.CreditIndex[ivfa.Credit].AddProfile(ivfa.ProfilePtr)
	}
}

func (ii *InvertIndex) UpdateActive(ivfb, ivfa InvertFeature) {
	if ivfb.ProfilePtr != nil {
		ii.ActiveIndex[ivfb.Active].RemoveProfile()
	}
	if ii.ActiveIndex[ivfa.Active] == nil {
		ii.ActiveIndex[ivfa.Active] = &ProfileSliceCount{}
	}
	if ivfa.ProfilePtr != nil {
		ii.ActiveIndex[ivfa.Active].AddProfile(ivfa.ProfilePtr)
	}
}

func (ii *InvertIndex) UpdateFaceCount(ivfb, ivfa InvertFeature) {
	if ivfb.ProfilePtr != nil {
		ii.FaceCountIndex[ivfb.FaceCount].RemoveProfile()
	}
	if ii.FaceCountIndex[ivfa.FaceCount] == nil {
		ii.FaceCountIndex[ivfa.FaceCount] = &ProfileSliceCount{}
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
	memory      	*Memory
	invertIndex 	*InvertIndex
	requestCache	*RequestCache
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
	fmt.Printf("init RequestCache")
	sp.requestCache.Init(sp.memory)
}



func MatchByInfo(profile *Profile, match []string) bool {
	time.Sleep(time.Duration(rand.Int63n(MatchByInfoCost)) * time.Nanosecond)
	for _, item := range match {
		if item == strconv.Itoa(int(profile.Match)) {
			return true
		}
	}
	return false
}

func (sp *SuperProfile) SearchByInfo(request *Request) ([]*Profile, int) {
	index := request.Index
	offset := request.Offset
	limit := request.Limit
	match := strings.Split(request.Match, ",")
	var result []*Profile
	sp.memory.RLock()
	for _, profile := range sp.memory.Profiles {
		if MatchByInfo(profile, match) {
			result = append(result, profile)
		}
	}
	sort.Sort(ProfileSlice{
		Profiles:  result,
		SortIndex: index,
	})
	sp.memory.RUnlock()
	if len(result) < offset {
		return nil, len(result)
	} else {
		if len(result) < offset + limit {
			return result[offset: ], len(result)
		} else {
			return result[offset: offset + limit], len(result)
		}
	}
}

type Request struct {
	Match string
	Index int
	Offset int
	Limit int
}


func (sp *SuperProfile) SearchByInfoIndex(request *Request) ([]*Profile, int) {
	defer func() {
		sp.invertIndex.Unlock()
	}()
	index := request.Index
	match := strings.Split(request.Match, ",")
	offset := request.Offset
	limit := request.Limit
	total := sp.requestCache.GetTotalAndUpdate(request.Match, sp.memory)
	if total < offset {
		return nil, total
	}
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
		profileListWapper := indexItem[key]
		profileList := profileListWapper.Profiles
		if !profileListWapper.RemoveFlag {
			for _, p := range profileList {
				if p != sp.memory.Profiles[p.ProfileID] {
					fmt.Printf("not equal to memorys, key = %d\n%v\n%v\n", key, p, sp.memory.Profiles[p.ProfileID])
				}
				if MatchByInfo(p, match) {
					results = append(results, p)
					if len(results) == offset + limit {
						return results[offset: offset + limit], total
					}
				}
			}
		} else {
			tmpSet := make(map[string]bool, len(profileList))
			sp.memory.RLock()
			result := make([]*Profile, 0, len(profileList))
			for _, profile := range profileList {
				switch index {
				case Credit:
					if profile.Credit != key {
						continue
					}
				case Active:
					if profile.Active != key {
						continue
					}
				case FaceCount:
					if profile.FaceCount != key {
						continue
					}
				default:
					if profile.Credit != key {
						continue
					}
				}
				tmpProfile := sp.memory.Profiles[profile.ProfileID]
				if tmpProfile != profile {
					continue
				}
				if tmpSet[profile.ProfileID] == true {
					fmt.Printf("already in tmpSet")
					continue
				} else {
					tmpSet[profile.ProfileID] = true
					result = append(result, tmpProfile)
					if MatchByInfo(tmpProfile, match) {
						results = append(results, tmpProfile)
						if len(results) == offset + limit {
							return results[offset: offset + limit], total
						}
					}
				}
			}
			indexItem[key].Reset(result)
		}
	}
	return results, len(results)
}

type CacheData struct {
	TotalMatched 	int
	LastRequestTime	time.Time
	IsDefault		bool
}

func (cd *CacheData) mm() {
	cd.TotalMatched--
}

func (cd *CacheData) pp() {
	cd.TotalMatched++
}

func (cd *CacheData) Update() {
	cd.LastRequestTime = time.Now()
}

type RequestCache struct {
	Requests	map[string]*CacheData
	sync.RWMutex
}

func (rc *RequestCache) Init(memory *Memory) {
	rc.Requests = make(map[string]*CacheData, 50)
	memory.RLock()
	rc.Requests["all"] = &CacheData{
		TotalMatched: len(memory.Profiles),
		LastRequestTime: time.Unix(2147483647, 0),
		IsDefault: true,
	}
	cnt := 0
	for _, p := range memory.Profiles {
		if p.Match < 10 {
			cnt ++
		}
	}
	rc.Requests["0,1,2,3,4,5,6,7,8,9"] = &CacheData{
		TotalMatched: cnt,
		LastRequestTime: time.Unix(2147483647, 0),
		IsDefault: true,
	}

	memory.RUnlock()
}

func (rc *RequestCache) MinusMinus(profile *Profile) {
	rc.Lock()
	for key := range rc.Requests {
		keyItems := strings.Split(key, ",")
		for _, item := range keyItems {
			if strconv.Itoa(int(profile.Match)) == item {
				rc.Requests[key].mm()
				break
			}
		}
	}
	rc.Unlock()
}

func (rc *RequestCache) PlusPlus(profile *Profile) {
	rc.Lock()
	for key := range rc.Requests {
		keyItems := strings.Split(key, ",")
		for _, item := range keyItems {
			if strconv.Itoa(int(profile.Match)) == item {
				rc.Requests[key].pp()
				break
			}
		}
	}
	rc.Unlock()
}

func (rc *RequestCache) GetTotalAndUpdate(req string, memory *Memory) int {
	rc.Lock()
	for key, value := range rc.Requests {
		if key == req {
			rc.Requests[req].Update()
			rc.Unlock()
			return value.TotalMatched
		}
	}
	rc.Unlock()

	// 新增一个 requestCache
	reqItemsList := strings.Split(req, ",")
	reqItemsSet := make(map[string]bool)
	for _, item := range reqItemsList {
		reqItemsSet[item] = true
	}
	memory.RLock()
	var cnt int
	for _, p := range memory.Profiles {
		if reqItemsSet[strconv.Itoa(int(p.Match))] == true {
			cnt++
		}
	}
	memory.RUnlock()
	rc.Lock()
	rc.Requests[req] = &CacheData{
		TotalMatched: cnt,
		LastRequestTime: time.Now(),
		IsDefault: false,
	}
	if len(rc.Requests) > 50 {
		var oldest = time.Now()
		var oldestKey = ""
		for key, value := range rc.Requests {
			if value.LastRequestTime.Before(oldest) {
				oldest = value.LastRequestTime
				oldestKey = key
			}
		}
		if oldestKey != "" {
			delete(rc.Requests, oldestKey)
		}
	}
	rc.Unlock()
	return cnt
}