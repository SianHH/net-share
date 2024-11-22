package file_storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileDataInterface interface {
	GetKey() string
}

// 文件存储
type FileStorage[T FileDataInterface] struct {
	dbFile  string        // 文件保存位置
	second  int           // 持久化时间 0表示
	data    []T           // 内存数据
	dataKey *sync.Map     // 唯一key标识
	lock    *sync.RWMutex // 读写锁
	isOpen  bool          // 是否运行
}

func NewFileStorage[T FileDataInterface](target T, dbFile string, second int) (*FileStorage[T], error) {
	result := FileStorage[T]{
		dbFile:  dbFile,
		second:  second,
		data:    []T{},
		dataKey: &sync.Map{},
		lock:    &sync.RWMutex{},
		isOpen:  true,
	}
	if err := result.load(); err != nil {
		return nil, err
	}
	go result.runSaveTask()
	return &result, nil
}

func (fs *FileStorage[T]) Create(data T) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	if !fs.isOpen {
		return errors.New("is not open")
	}
	if _, ok := fs.dataKey.Load(data.GetKey()); ok {
		return errors.New("duplicate key")
	}
	fs.dataKey.Store(data.GetKey(), true)
	fs.data = append(fs.data, data)
	return fs.saveNoLock()
}

func (fs *FileStorage[T]) Delete(key string) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	if !fs.isOpen {
		return errors.New("is not open")
	}
	if _, ok := fs.dataKey.Load(key); !ok {
		return errors.New("key is not exist")
	}
	fs.dataKey.Delete(key)
	var newData []T
	for _, item := range fs.data {
		if item.GetKey() != key {
			newData = append(newData, item)
		}
	}
	fs.data = newData
	return fs.saveNoLock()
}

func (fs *FileStorage[T]) Update(data T) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	if !fs.isOpen {
		return errors.New("is not open")
	}
	if _, ok := fs.dataKey.Load(data.GetKey()); ok {
		var newData []T
		for _, item := range fs.data {
			if item.GetKey() != data.GetKey() {
				newData = append(newData, item)
			} else {
				newData = append(newData, data)
			}
		}
		fs.data = newData
	} else {
		fs.dataKey.Store(data.GetKey(), true)
		fs.data = append(fs.data, data)
	}
	return fs.saveNoLock()
}

func (fs *FileStorage[T]) Query(key string) (result T, err error) {
	fs.lock.RLock()
	defer fs.lock.RUnlock()
	if !fs.isOpen {
		return result, errors.New("is not open")
	}
	if _, ok := fs.dataKey.Load(key); !ok {
		return result, errors.New("key is not exist")
	}
	for _, item := range fs.data {
		if item.GetKey() == key {
			return item, nil
		}
	}
	return result, errors.New("data is not exist")
}

func (fs *FileStorage[T]) QueryFilter(filter ...func(T) bool) (result T, err error) {
	fs.lock.RLock()
	defer fs.lock.RUnlock()
	if !fs.isOpen {
		return result, errors.New("is not open")
	}
	for _, item := range fs.data {
		var success = len(filter)
		for _, f := range filter {
			if f(item) {
				success--
			}
		}
		if success == 0 {
			return item, nil
		}
	}
	return result, errors.New("not qualified data")
}

func (fs *FileStorage[T]) QueryAllFilter(filter ...func(T) bool) (result []T) {
	fs.lock.RLock()
	defer fs.lock.RUnlock()
	if !fs.isOpen {
		return nil
	}
	for _, item := range fs.data {
		var success = len(filter)
		for _, f := range filter {
			if f(item) {
				success--
			}
		}
		if success == 0 {
			result = append(result, item)
		}
	}
	return result
}

func (fs *FileStorage[T]) QueryAll() []T {
	fs.lock.RLock()
	defer fs.lock.RUnlock()
	if !fs.isOpen {
		return nil
	}
	return fs.data
}

func (fs *FileStorage[T]) Close() {
	if !fs.isOpen {
		return
	}
	_ = fs.save()
	fs.isOpen = false
}

func (fs *FileStorage[T]) runSaveTask() {
	for {
		// 判断是否关闭
		if !fs.isOpen || fs.second <= 0 {
			return
		}
		time.Sleep(time.Second * time.Duration(fs.second))
		if err := fs.save(); err != nil {
			panic("timer save file err:" + err.Error())
		}
	}
}

func (fs *FileStorage[T]) load() error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	if err := os.MkdirAll(filepath.Dir(fs.dbFile), os.ModeDir); err != nil {
		return err
	}
	file, err := os.OpenFile(fs.dbFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	dataBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(dataBytes, &fs.data)
	for _, item := range fs.data {
		fs.dataKey.Store(item.GetKey(), true)
	}
	return nil
}

func (fs *FileStorage[T]) save() error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	file, err := os.OpenFile(fs.dbFile, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	marshal, err := json.Marshal(fs.data)
	if err != nil {
		return err
	}
	if _, err = file.Write(marshal); err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage[T]) saveNoLock() error {
	// 如果开启定时持久化，则不立即执行保存
	if fs.second > 0 {
		return nil
	}
	file, err := os.OpenFile(fs.dbFile, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	marshal, err := json.Marshal(fs.data)
	if err != nil {
		return err
	}
	if _, err = file.Write(marshal); err != nil {
		return err
	}
	return nil
}
