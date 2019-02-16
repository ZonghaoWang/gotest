package convertindex

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)



func ChangeProfile(profile *Profile) *Profile {
	if rand.Int63n(10) >= 4 {
		if profile == nil {
			return &Profile{
				ProfileID: strconv.Itoa(leng + rand.Intn(65535)),
				Active: rand.Int63n(2 << MaxActiveBitShift - 1),
				Credit: rand.Int63n(2 << MaxCreditBitShift - 1),
				FaceCount: rand.Int63n(2 << MaxFaceCountBitShift - 1),
				Match: rand.Int63n(MatchResultProp),
			}
		} else {
			profile.Active = rand.Int63n(2 << MaxActiveBitShift - 1)
			profile.Credit = rand.Int63n(2 << MaxCreditBitShift - 1)
			profile.FaceCount = rand.Int63n(2 << MaxFaceCountBitShift - 1)
			profile.Match = rand.Int63n(MatchResultProp)
			return profile
		}
	} else {
		return nil
	}

}

func GetmockProfile() *Profile {
	if rand.Int63n(10) >= 4 {
		return &Profile{
			ProfileID: strconv.Itoa(leng + rand.Intn(65535)),
			Active: rand.Int63n(2 << MaxActiveBitShift - 1),
			Credit: rand.Int63n(2 << MaxCreditBitShift - 1),
			FaceCount: rand.Int63n(2 << MaxFaceCountBitShift - 1),
			Match: rand.Int63n(MatchResultProp),
		}
	} else {
		return nil
	}
}

func (sp *SuperProfile) CURDProfiles()  {
	var flag int
	for {
		flag++
		if flag > 100000 {
			break
		}
		profile := GetmockProfile()
		if profile != nil {
			if sp.memory.Profiles[profile.ProfileID] == nil {
				profile = nil
			} else {
				profile = sp.memory.Profiles[profile.ProfileID]
			}
		}
		if profile != nil {
			sp.memory.Lock()
			delete(sp.memory.Profiles, profile.ProfileID)
			sp.memory.Unlock()
		}
		//fmt.Printf("before profiel is %v\n", profile)
		cacheBefor := CacheFeature(profile)
		afterProfile := ChangeProfile(profile)
		//fmt.Printf("after profile is %v\n", afterProfile)
		if afterProfile != nil {
			sp.memory.Lock()
			sp.memory.Profiles[afterProfile.ProfileID] = afterProfile
			sp.memory.Unlock()
		}
		cacheAfter := CacheFeature(afterProfile)
		//fmt.Printf("before %v\n", cacheBefor)
		//fmt.Printf("after %v\n", cacheAfter)
		sp.invertIndex.Lock()
		sp.invertIndex.Update(cacheBefor, cacheAfter)
		sp.invertIndex.Unlock()
		time.Sleep(time.Microsecond)
	}
}


func BenchmarkSearch1(b *testing.B) {
	sp := SuperProfile{
		memory:      &Memory{},
		invertIndex: &InvertIndex{},
	}
	sp.Init()
	//go sp.CURDProfiles()
	//time.Sleep(time.Second * 4)
	fmt.Println("length of profiles is ", len(sp.memory.Profiles))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Println("length of result is ,", len(sp.SearchByInfo(Credit)))
	}
}

func BenchmarkSearch2(b *testing.B) {
	sp := &SuperProfile{
		memory:      &Memory{},
		invertIndex: &InvertIndex{},
	}
	sp.Init()
	//go sp.CURDProfiles()
	//time.Sleep(time.Second * 4)
	fmt.Println("length of profiles is ", len(sp.memory.Profiles))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Println("length of result is ,", len(sp.SearchByInfoIndex(Credit)))
	}
}

func TestTwoSearch(t *testing.T)  {
	sp := &SuperProfile{
		memory:      &Memory{},
		invertIndex: &InvertIndex{},
	}
	sp.Init()
	go sp.CURDProfiles()
	time.Sleep(time.Second * 2)
	fmt.Println("length of memory is ", len(sp.memory.Profiles))
	profiles0 := sp.SearchByInfo(Credit)
	fmt.Println("length of profiles0 is ", len(profiles0))
	idSet0 := make(map[string]bool)
	for _, p := range profiles0 {
		idSet0[p.ProfileID] = true
	}
	fmt.Println("len of idset0 is ", len(idSet0))
	profiles := sp.SearchByInfoIndex(Credit)
	idSet := make(map[string]bool)
	for _, p := range profiles {
		if idSet0[p.ProfileID] == false {
			fmt.Printf("not in idset0 %v\n", p)
			fmt.Printf("%v\n", sp.memory.Profiles[p.ProfileID])
		}

		if sp.memory.Profiles[p.ProfileID] == nil {
			fmt.Printf("not in memory %v\n", *p)
		}

		if p.Match >= 10 {
			fmt.Printf("not match %v\n", *p)
		}

		if idSet[p.ProfileID] == true {
			fmt.Printf("colle %v\n", *p)
		}
		idSet[p.ProfileID] = true
	}


	for pid, _ := range sp.memory.Profiles {
		if idSet[pid] == false {
			fmt.Printf("not in memory %s\n", pid)
		}
	}

	fmt.Println("length of set is ", len(profiles))
	fmt.Println("length of set is ", len(idSet))
}

func BenchmarkAppendSlice(b *testing.B) {
	rst := make([]*Profile, 0, 1000)
	rst2 := make([]*Profile, 1000)
	var rrst []*Profile
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rrst = append(rst, rst2...)
	}
	fmt.Println(len(rrst))
}
