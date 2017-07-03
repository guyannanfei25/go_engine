package consumer

import (
    "fmt"
    "github.com/guyannanfei25/go-logger/logger"
    sj  "github.com/guyannanfei25/go-simplejson"
)

// use MultiHandlerDispatcher simple, no need to consider mulit goroutine.
type Echoer struct {
    num     int
    id      int
}

func (e *Echoer) Init(conf *sj.Json, id int) error {
    e.num = 0
    e.id  = id
    return nil
}

func (e *Echoer) Process(item interface{}) (interface{}, error) {
    e.num += 1
    rItem, ok := item.(*Context)
    if !ok {
        logger.Errorf("Unexpected type[%T]\n", item)
        return nil, fmt.Errorf("Unexpected type[%T]", item)
    }

    logger.Infof("Get raw log len[%d] Bytes\n", len(rItem.Raw))

    return item, nil
}

func (e *Echoer) Tick() {
    logger.Infof("now %dth echoer get [%d] items\n", e.id, e.num)
}

func (e *Echoer) Close() {
    logger.Infof("%dth echoer close, finish process [%d] items\n", e.id, e.num)
}
