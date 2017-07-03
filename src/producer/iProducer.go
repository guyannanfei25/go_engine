package producer

import (
    "errors"
    sj  "github.com/guyannanfei25/go-simplejson"
)

var (
    ErrReadTimeout = errors.New("producer: read timeout")
)

// Producer is an interface general data.
// It is the data provider's duty to provide id in order to remove duplicate.
type Producer interface {
    Init(ctx *sj.Json) error
    GetMsg() ([]byte, error)
    Close() error
}

func ProducerFactory(pType string) Producer {
    switch pType {
    case "nsq":
        return new(NsqProducer)
    default:
        return nil
    }
}
