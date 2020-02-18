package services

import (
	"github.com/astaxie/beego/logs"
	"sync"
	"time"
)

//TaskFun 执行任务的方法
type TaskFun func()

//Task 任务
type Task struct {
	lastTime time.Time
	d time.Duration
	f TaskFun
	name string
	doing bool //是否正在执行
	lock sync.Mutex
}

//execSingle 单任务异步执行任务
func (this *Task) execSingle() {
	logs.Debug("task:%s 定时器监测到需要执行,lastTime:%s", this.name, this.lastTime.Format("2006-01-02 15:04:05"))

	//排除重复执行
	this.lock.Lock()
	if this.doing {
		logs.Debug("task:%s 存在进程正在执行，本次执行中断", this.name)
		this.lock.Unlock()
		return
	}
	this.doing = true
	this.lock.Unlock()

	this.lastTime = time.Now()

	//异步执行任务
	go func() {
		//定义异常记录
		defer func() {
			if err := recover(); err != nil {
				logs.Error("task:%s 异常:%s", this.name, err)
			}
			this.lock.Lock()
			this.doing = false
			this.lock.Unlock()
		}()
		//执行任务
		logs.Debug("task:%s 开始执行任务", this.name)
		this.f()
	}()

}

//TaskService 任务服务
type TaskService struct {
	taskList []*Task
}

//AddTask 创建任务
func (this *TaskService) AddTask(name string, d time.Duration, f TaskFun) {
	t := &Task{
		d:d,
		f:f,
		name:name,
		doing:false,
	}
	this.taskList = append(this.taskList, t)
}

//TaskSingleStart 开始单任务执行启动
func (this *TaskService)TaskSingleStart() {
	logs.Info("启动单任务处理定时任务服务...")
	for _,task := range this.taskList {
		task.execSingle()
	}
	//定义秒级定时器
	t := time.NewTicker(time.Second)

	for _ = range t.C {
		for _,task := range this.taskList {
			if time.Now().Sub(task.lastTime) >= task.d {
				task.execSingle()
			}
		}
	}
}