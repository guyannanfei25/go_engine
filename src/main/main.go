package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "sync"
    "time"

    "producer"
    "consumer"

    // "github.com/guyannanfei25/flowprocess"
    "github.com/guyannanfei25/go_glue"
    "github.com/guyannanfei25/go-logger/logger"
    sj  "github.com/guyannanfei25/go-simplejson"

)

var conf = flag.String("f", "etc/example.json", "json conf path")

func main() {
    flag.Parse()
    ctx, err := sj.NewFromFile(*conf)
    if err != nil {
        fmt.Fprintf(os.Stderr, "parse conf[%s] to json err[%s]\n", *conf, err)
        panic(fmt.Sprintf("parse conf[%s] to json err[%s]", *conf, err))
    }

    if err := glue.Init(ctx); err != nil {
        fmt.Fprintf(os.Stderr, "glue init err[%s]\n", err)
        panic(fmt.Sprintf("glue init err[%s]", err))
    }

    // init producer
    p := producer.ProducerFactory(ctx.Get("producer").Get("type").MustString("nsq"))
    if p == nil {
        panic("new producer err!")
    }

    if err := p.Init(ctx.Get("producer")); err != nil {
        logger.Errorf("new producer err[%s]\n", err)
        panic("new producer err")
    }
    defer p.Close()

    // init framework
    framework := InitFramework(ctx.Get("consumers"))
    defer framework.Close()
    framework.Start()

    // init signal
    notifySig := make(chan os.Signal, 1)
    signal.Notify(notifySig, syscall.SIGINT, syscall.SIGTERM)

    notify := make(chan struct{}, 1)

    wg := new(sync.WaitGroup)
    wg.Add(3)

    // notify
    go func() {
        <- notifySig
        logger.Debugf("Receive sig\n")
        close(notify)
        wg.Done()
    } ()

    // init tick
    go func() {
        defer wg.Done()
        tickTime := ctx.Get("consumers").Get("tick_s").MustInt(20)
        ticker := time.NewTicker(time.Duration(tickTime) * time.Second)
TICK_LOOP:
        for {
            select {
            case <- ticker.C:
                logger.Debugf("Begin tick")
                framework.Tick()
            case <- notify:
                logger.Debugf("tick got exit notify")
                ticker.Stop()
                break TICK_LOOP
            }
        }
    }()

    go func() {
        defer wg.Done()
MAIN_LOOP:
        for {
            select {
            case <- notify:
                logger.Infof("Get sig to exit, no exiting\n")
                break MAIN_LOOP
            default:
                msg, err := p.GetMsg()
                if err != nil {
                    logger.Errorf("Get msg err[%s] from src providuer\n", err)
                    continue
                }
                c := new(consumer.Context)
                c.Raw = msg
                framework.Dispatch(c)
            }
        }
    }()

    wg.Wait()
    logger.Infof("main program exit.\n")
}
