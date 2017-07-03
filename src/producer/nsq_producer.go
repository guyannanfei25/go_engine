package producer

import (
    "time"
    "fmt"
    "github.com/nsqio/go-nsq"
    "github.com/guyannanfei25/go-logger/logger"
    sj  "github.com/guyannanfei25/go-simplejson"
)

type NsqProducer struct {
    // TODO: multi nsq consumer
    consumer        *nsq.Consumer
    config          *nsq.Config
    timeout         int // GetMsg read timeout second
    msgChan         chan *nsq.Message
    topic           string
    channel         string
}

func (p *NsqProducer) Init(ctx *sj.Json) error {
    if ctx == nil {
        return fmt.Errorf("init nsq prodcuer err, no json content")
    }

    p.timeout    = ctx.Get("read_timeout_s").MustInt(10)
    bufSize     := ctx.Get("nsq_consumer_conf").Get("buf_size").MustInt(0)
    if bufSize < 0 {
        bufSize = 0
    }
    p.msgChan    = make(chan *nsq.Message, bufSize)
    nsqLookupds := ctx.Get("nsq_consumer_conf").Get("lookupds").MustStringArray()
    if len(nsqLookupds) < 1 {
        return fmt.Errorf("nsq prodcuer conf get no nsqlookupds addr")
    }

    p.topic      = ctx.Get("nsq_consumer_conf").Get("topic").MustString("")
    p.channel    = ctx.Get("nsq_consumer_conf").Get("channel").MustString("")
    if p.topic == "" || p.channel == "" {
        return fmt.Errorf("init nsq producer conf get no topic or channel")
    }

    p.config    = nsq.NewConfig()
    p.config.MaxInFlight     = ctx.Get("nsq_consumer_conf").Get("max_in_flight").MustInt(100)
    p.config.DialTimeout     = time.Duration(ctx.Get("nsq_consumer_conf").Get("dail_timeout_s").MustInt(10)) * time.Second
    p.config.ReadTimeout     = time.Duration(ctx.Get("nsq_consumer_conf").Get("read_timeout_s").MustInt(60)) * time.Second
    p.config.WriteTimeout    = time.Duration(ctx.Get("nsq_consumer_conf").Get("write_timeout_s").MustInt(30)) * time.Second
    p.config.LookupdPollInterval  = time.Duration(ctx.Get("nsq_consumer_conf").Get("lookupd_poll_interval_s").MustInt(20)) * time.Second

    var err error
    p.consumer, err = nsq.NewConsumer(p.topic, p.channel, p.config)
    if err != nil {
        return fmt.Errorf("NsqProducer new nsq consumer err[%s]", err)
    }

    p.consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
        p.msgChan <- msg
        // TODO: not use default???
        return nil
    }))

    err = p.consumer.ConnectToNSQLookupds(nsqLookupds)
    if err != nil {
        return fmt.Errorf("NsqProducer connect to nsqlookupds[%v] err[%s]", nsqLookupds, err)
    }

    logger.Debugf("NsqProducer init success, nsq consumer read timeout[%d]s\n", p.timeout)

    return nil
}

func (p *NsqProducer) GetMsg() ([]byte, error) {
    select {
    case msg := <- p.msgChan:
        return msg.Body, nil
    case <- time.After(time.Duration(p.timeout) * time.Second):
        return nil, ErrReadTimeout
    }
}

func (p *NsqProducer) Close() error {
    logger.Debugf("NsqProducer start exiting\n")
    p.consumer.Stop()
    <- p.consumer.StopChan
    close(p.msgChan)
    logger.Debugf("NsqProducer exit success\n")
    return nil
}
