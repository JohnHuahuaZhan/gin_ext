package session

import (
	"sync"
	"time"
)

type Storager interface {
	//初始化一个session，id根据需要生成后传入
	Init(sid string, timeout time.Duration) (Sessioner, error)
	//根据sid，获得当前session
	Get(sid string) (Sessioner, error)
	//销毁session
	Destroy(sid string) error
	//回收
	GC()
}

type MemoryStorage struct {
	lock             sync.RWMutex                 //一把读写锁
	data             map[string]Sessioner //数据
}
//实例化一个内存实现
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]Sessioner, 100),
	}
}
func (memoryStorage *MemoryStorage) Init(sid string, timeout time.Duration) (Sessioner, error){

	if len(sid) == 0{
		return nil,ErrorSessionIdNotBeEmpty
	}
	s := NewCommonSession()
	s.sid = sid
	s.timeout = timeout
	memoryStorage.lock.Lock()
	defer memoryStorage.lock.Unlock()
	memoryStorage.data[sid] = s
	return s, nil
}

//更新access time 是manager的事情
func (memoryStorage *MemoryStorage)Get(sid string) (Sessioner, error){
	memoryStorage.lock.RLock()
	defer memoryStorage.lock.RUnlock()
	if s, ok := memoryStorage.data[sid]; ok{
		return s, nil
	}else {
		return  nil, ErrorSessionNotFound
	}
}

func (memoryStorage *MemoryStorage) Destroy(sid string) error{
	memoryStorage.lock.Lock()
	defer memoryStorage.lock.Unlock()
	delete(memoryStorage.data, sid)
	return nil
}
func (memoryStorage *MemoryStorage) GC() {
	memoryStorage.lock.Lock()
	defer memoryStorage.lock.Unlock()
	for sessionID, session := range memoryStorage.data {
		//删除超过时限的session
		if session.Access().Add(session.Timeout()).Before(time.Now().UTC()) {
			delete(memoryStorage.data, sessionID)
		}
	}
}