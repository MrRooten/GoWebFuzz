package fuzz

import (
	"bufio"
	"gowebfuzz/library/fuzz/fuzzlib"
	"gowebfuzz/library/network/gwfhttp"
	"gowebfuzz/library/utils"
	"gowebfuzz/library/utils/goroutinelib"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type FuzzType int

const (
	Lua FuzzType = iota
	Golang FuzzType = 1
	Json FuzzType = 2
	Yaml FuzzType = 3
)

type RequestNode struct {
	Next *RequestNode
	Last *RequestNode
	Request *gwfhttp.RequestPacket
}

type ReuqestQueue struct {
	lock sync.Mutex
	Head *RequestNode
	End *RequestNode
	Len int
}

func (queue *ReuqestQueue) Push(req *gwfhttp.RequestPacket) {
	queue.lock.Lock()
	defer queue.lock.Unlock()
	if queue.Head == nil {
		reqNode := RequestNode{
			Next: nil,
			Last: nil,
			Request: req,
		}
		queue.Head = &reqNode
		queue.End = &reqNode
		queue.Len = 1
		return
	}

	reqNode := RequestNode{
		Next: queue.End,
		Last: nil,
		Request: req,
	}
	reqNode.Next.Last = &reqNode
	queue.End = &reqNode
	queue.Len += 1
}

func (queue *ReuqestQueue) Take() *gwfhttp.RequestPacket {
	queue.lock.Lock()
	defer queue.lock.Unlock()
	if queue.Len == 0 {
		return nil
	}
	reqNode := queue.Head
	queue.Head = queue.Head.Last
	if queue.Len == 1 {
		queue.End = nil
	}
	queue.Len -= 1
	return reqNode.Request
}
var logger utils.Logger

type Fuzz struct {
	pool goroutinelib.GoroutinePool
	waitingQueue ReuqestQueue
	IsInit bool
	RoutineSize int
	ModuleRoutineSize int
	queueChan chan *gwfhttp.RequestPacket
	RunningFuzzType FuzzType
	FuzzPath string
	FuzzFile string
}

func GetPathModuleType(path string) FuzzType {
	if strings.HasPrefix(path,"lua://") {
		return Lua
	} else if strings.HasPrefix(path,"go://") {
		return Golang
	} else if strings.HasPrefix(path,"json://") {
		return Json
	} else if strings.HasPrefix(path,"yaml://") {
		return Yaml
	}

	return -1
}

func GetPathValue(path string) string {
	return ""
}

func SplitFuzzPath(paths string) []string {
	res := []string{}
	res = strings.Split(paths,",")
	return res
}

func ReadFuzzFromFile(file string) []string {
	var res []string
	logger := utils.Logger{}
	fi,err := os.Open(file)
	if err != nil {
		logger.Warningf("The file open error: %s",err.Error())
		return res
	}
	defer func(fi *os.File) {
		err := fi.Close()
		if err != nil {
			logger.Warningf("The file close error: %s",err.Error())
		}
	}(fi)

	br := bufio.NewReader(fi)
	for {
		line,_,c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if len(strings.TrimSpace(string(line))) != 0 {
			res = append(res,string(line))
		}
	}
	return res
}
//modulesRun(modules *[]fuzzlib.Module,reqp *gwfhttp.RequestPacket)
//Every request has a goroutine pool that size is fuzz.ModuleRoutineSize
func (fuzz *Fuzz)modulesRun(vargs []interface{}) {
	modules := vargs[0].(*[]fuzzlib.Module)
	reqp := vargs[1].(*gwfhttp.RequestPacket)

	pool := goroutinelib.GoroutinePool{}
	pool.NewGoroutinePool(fuzz.ModuleRoutineSize)
	for _,module := range *modules {
		var moduleVargs []interface{}
		moduleVargs = append(moduleVargs,module)
		moduleVargs = append(moduleVargs,reqp)
		pool.RunTask(func(vargs []interface{}) {
			module := vargs[0].(fuzzlib.Module)
			reqp := vargs[1].(*gwfhttp.RequestPacket)
			err := module.Run(reqp)
			if err != nil {
				logger.Error(err.Error())
				return
			}
		},moduleVargs)
	}
}
func (fuzz *Fuzz) StartFuzzing() {
	var modulePaths []string
	var modules []fuzzlib.Module

	if len(fuzz.FuzzPath) != 0 {
		modulePaths = SplitFuzzPath(fuzz.FuzzPath)
	} else if len(fuzz.FuzzFile) != 0 {
		modulePaths = ReadFuzzFromFile(fuzz.FuzzFile)
	}
	//To get every module type
	for _,modulePath := range modulePaths {
		var moduleType FuzzType
		moduleType = GetPathModuleType(modulePath)
		var module fuzzlib.Module
		if moduleType == Lua {
			module = fuzzlib.NewLuaModule(modulePath)
		} else if moduleType == Golang {
			module = fuzzlib.NewGolangModule(modulePath)
		} else if moduleType == Json {
			module = fuzzlib.NewJsonModule(modulePath)
		} else if moduleType == Yaml {
			module = fuzzlib.NewYamlModule(modulePath)
		} else {
			logger.Warningf("The module %s type is not existed",modulePath)
			continue
		}
		modules = append(modules,module)
	}

	//The Global queue has a goroutine pool that size is fuzz.routineSize
	for true {
		reqp := fuzz.waitingQueue.Take()
		//fmt.Println(fuzz.waitingQueue.Len)
		if reqp == nil {
			continue
		}
		var vargs []interface{}
		vargs = append(vargs,&modules)
		vargs = append(vargs,reqp)
		fuzz.pool.RunTask(fuzz.modulesRun,vargs)
	}
}

//Passive fuzz only match uri body and POST body
func (fuzz *Fuzz)AppendQueue(requeset *http.Request) error {
	reqp := gwfhttp.RequestPacket{}
	err := reqp.ReadFromHTTPRequest(requeset)
	if err != nil {
		return err
	}

	//Save the req into queue
	fuzz.waitingQueue.Push(&reqp)
	return nil
}

func (fuzz *Fuzz) Init(routineSize int,fuzzPath string,fuzzFile string,moduleRoutineSize int) {
	fuzz.pool.NewGoroutinePool(routineSize)
	fuzz.queueChan = make(chan *gwfhttp.RequestPacket,routineSize)
	fuzz.FuzzPath = fuzzPath
	fuzz.FuzzFile = fuzzFile
	fuzz.ModuleRoutineSize = moduleRoutineSize
}
func (fuzz *Fuzz)Drive(request *http.Request) {
	err := fuzz.AppendQueue(request)
	if err != nil {

	}
}