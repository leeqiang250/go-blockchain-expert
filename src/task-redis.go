package src

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type TaskRedis struct {
	client         *redis.Client
	script         string
	scriptSha1     string
	interval       time.Duration
	keyUniqueName  string
	nodeUniqueName string
	task           func()
}

func NewTaskRedis(client *redis.Client, script string, interval time.Duration, keyUniqueName string, nodeUniqueName string, task func()) *TaskRedis {
	return &TaskRedis{
		client:         client,
		script:         script,
		interval:       interval,
		keyUniqueName:  keyUniqueName,
		nodeUniqueName: nodeUniqueName,
		task:           task,
	}
}

func (this *TaskRedis) Run() {
	this.init()

	var ticker = time.NewTicker(this.interval)
	for {
		<-ticker.C

		if this.tryLock() {
			this.task()
		}
	}
}

func (this *TaskRedis) init() {
	var script, err = this.client.ScriptLoad(context.Background(), this.script).Result()
	if nil != err {
		//TODO
		//异常处理
	}

	this.scriptSha1 = script
}

func (this *TaskRedis) tryLock() bool {
	var result, err = this.client.EvalSha(
		context.Background(),
		this.scriptSha1,
		[]string{this.keyUniqueName},
		this.nodeUniqueName,
		uint8(this.interval.Seconds())+3,
	).Result()

	if nil != err {
		//TODO
		//异常处理
		return false
	}

	fmt.Println(this.keyUniqueName, this.nodeUniqueName, result)

	return "1" == fmt.Sprintf("%v", result)
}

//if (redis.call('get', KEYS[1]) == ARGV[1])
//then
//redis.call('expire', KEYS[1], ARGV[2])
//return 1
//else
//if (1 == redis.call('setnx', KEYS[1], ARGV[1]))
//then
//redis.call('expire', KEYS[1], ARGV[2])
//return 1
//else
//return 0
//end
//return 0
//end
