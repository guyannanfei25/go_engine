package consumer

import (
    "github.com/guyannanfei25/flowprocess"
)

// 工厂方法模式
func CreatorFactory(ctype string) flowprocess.HandlerCreator {
    switch ctype {
    case "Echoer":
        return func() flowprocess.Handler {
            return new(Echoer)
        }
    default:
        return nil
    }
}
