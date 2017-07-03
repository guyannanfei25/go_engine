package main

import (
    "consumer"
    "github.com/guyannanfei25/flowprocess"
    "github.com/guyannanfei25/go-logger/logger"
    sj  "github.com/guyannanfei25/go-simplejson"
)

type FrameWork struct {
    dispatchers      map[string]flowprocess.MultiHandlerDispatcher
    IsInit           bool
    IsRunning        bool
    ctx              *sj.Json
}

func (f *FrameWork) Init(ctx *sj.Json) error {
    f.IsInit = false
    f.IsRunning = false
    f.ctx = ctx
    f.dispatchers = make(map[string]flowprocess.MultiHandlerDispatcher, 10)
    framework := new(flowprocess.DefaultMultiHandlerDispatcher)
    err := framework.Init(ctx.Get("framework"))
    if err != nil {
        logger.Errorf("Init framework err[%s]\n", err)
        return err
    }
    logger.Debugf("init framework dispatch\n")
    f.dispatchers["framework"] = framework

    downDispatchers := f.getDownDispatchers(ctx.Get("framework").Get("down_consumers").MustStringArray([]string{}))
    for _, d := range downDispatchers {
        framework.DownRegister(d)
    }
    f.IsInit = true

    logger.Debugf("framework init success")
    return nil
}

func (f *FrameWork) getDownDispatchers(dispatchers []string) []flowprocess.MultiHandlerDispatcher {
    logger.Debugf("getDownDispatchers dispatchers[%#v]\n", dispatchers)
    var ret = make([]flowprocess.MultiHandlerDispatcher, 0, len(dispatchers))
    if len(dispatchers) == 0 {
        return ret
    }

    // FIXME: topolgy, forever loop
    for _, dType := range dispatchers {
        var disp flowprocess.MultiHandlerDispatcher
        if tmpd, ok := f.dispatchers[dType]; ok {
            disp = tmpd
        } else {
            disp = new(flowprocess.DefaultMultiHandlerDispatcher)
            f.dispatchers[dType] = disp
            if err := disp.Init(f.ctx.Get(dType)); err != nil {
                logger.Errorf("init dispatcher type[%s] err[%s]\n", dType, err)
                panic(err)
            }
            logger.Debugf("Init type[%s] dispatcher\n", dType)
            disp.RegisterHandlerCreator(consumer.CreatorFactory(dType))

            downDis := f.getDownDispatchers(f.ctx.Get(dType).Get("down_consumers").MustStringArray([]string{}))
            for _, d := range downDis {
                disp.DownRegister(d)
            }
        }

        ret = append(ret, disp)
    }

    return ret
}

func (f *FrameWork) Start() {
    f.dispatchers["framework"].Start()
    f.IsRunning = true
}

func (f *FrameWork) Dispatch(item interface{}) {
    f.dispatchers["framework"].Dispatch(item)
}

func (f *FrameWork) Tick() {
    f.dispatchers["framework"].Tick()
}

func (f *FrameWork) Close() {
    f.dispatchers["framework"].Close()
    f.IsInit = false
    f.IsRunning = false
}

// ctx is consumer key content in json
// func InitFramework(ctx *sj.Json) flowprocess.MultiHandlerDispatcher {
    // dispatchers := make(map[string]flowprocess.MultiHandlerDispatcher)
    // framework := new(flowprocess.DefaultMultiHandlerDispatcher)
    // framework.Init(ctx.Get("framework"))
    // dispatchers["framework"] = framework

    // nowDispatcher := framework
    // var downDispatchers []string
    
    // for {

    // }
// }

var DefaultFrameWork FrameWork

func InitFramework(ctx *sj.Json) *FrameWork {
    if err := DefaultFrameWork.Init(ctx); err != nil {
        panic(err)
    }
    return &DefaultFrameWork
}
