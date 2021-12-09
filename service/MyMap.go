package service

import (
	"hash/crc32"
	"sync"
)

var GlobalMyMap *MyMap

type MyMap struct {
	sliceMap           []map[string]*UserInfo
	cpuNum             uint32
	sliceLock []*sync.RWMutex
}

type UserInfo struct {
	Uid    int
	UserID string
}

type SearchChannel struct {
	getChannel    chan *GetInfo
	resultChannel chan *UserInfo
	setChannel    chan *UserInfo
	myMap         map[string]*UserInfo
}

type SetInfo struct {
	sliceKey int
	*UserInfo
}

func (s *SetInfo) setSliceKey(cpuNum uint32) {
	key := crc32.ChecksumIEEE([]byte(s.UserID))
	s.sliceKey = int(key % cpuNum)
}

type GetInfo struct {
	sliceKey      int
	sliceMapKey   string
	ResultChannel chan *UserInfo
}

func (g *GetInfo) setSliceKey(cpuNum uint32) {
	key := crc32.ChecksumIEEE([]byte(g.sliceMapKey))
	g.sliceKey = int(key % cpuNum)
}

func init() {
	GlobalMyMap = new(MyMap)
	GlobalMyMap.cpuNum = uint32(50)
	//fmt.Println("查看一下cpu ", GlobalMyMap.cpuNum)
	GlobalMyMap.sliceMap = make([]map[string]*UserInfo, GlobalMyMap.cpuNum)
	GlobalMyMap.sliceLock = make([]*sync.RWMutex, GlobalMyMap.cpuNum)
	for i := 0; i < int(GlobalMyMap.cpuNum); i++ {
		GlobalMyMap.sliceMap[i] = make(map[string]*UserInfo)
		GlobalMyMap.sliceLock[i] = new(sync.RWMutex)
	}
	GlobalNormalMap = new(NormalMap)
	GlobalNormalMap.DB = make(map[string]*UserInfo)
	GlobalNormalMap.RWMutex = sync.RWMutex{}
}

func (m *MyMap) Set(userInfo *UserInfo) {
	index32 := crc32.ChecksumIEEE([]byte(userInfo.UserID))
	index := int(index32 % m.cpuNum)
	m.sliceLock[index].Lock()
	m.sliceMap[index][userInfo.UserID] = userInfo
	m.sliceLock[index].Unlock()
}

func (m *MyMap) Search(UserID string) (u *UserInfo){

	index32 := crc32.ChecksumIEEE([]byte(UserID))
	index := int(index32 % m.cpuNum)
	m.sliceLock[index].RLock()
	defer m.sliceLock[index].RUnlock()
	u = m.sliceMap[index][UserID]

	return
}


type NormalMap struct {
	sync.RWMutex
	DB map[string]*UserInfo
}

var GlobalNormalMap *NormalMap

func (n *NormalMap)Set(userInfo *UserInfo)  {
	n.Lock()
	defer n.Unlock()
	n.DB[userInfo.UserID] = userInfo
}

func (n *NormalMap)Get(userID string) (*UserInfo) {
	n.RLock()
	defer n.RUnlock()

	return n.DB[userID]
}