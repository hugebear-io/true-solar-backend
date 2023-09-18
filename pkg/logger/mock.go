package logger

type loggerMock struct{}

func NewLoggerMock() Logger {
	return &loggerMock{}
}

func (*loggerMock) Debug(string)                  {}
func (*loggerMock) Info(string)                   {}
func (*loggerMock) Warn(string)                   {}
func (*loggerMock) Error(error)                   {}
func (*loggerMock) Panic(error)                   {}
func (*loggerMock) Fatal(error)                   {}
func (*loggerMock) Debugf(string, ...interface{}) {}
func (*loggerMock) Infof(string, ...interface{})  {}
func (*loggerMock) Warnf(string, ...interface{})  {}
func (*loggerMock) Errorf(string, ...interface{}) {}
func (*loggerMock) Panicf(string, ...interface{}) {}
func (*loggerMock) Fatalf(string, ...interface{}) {}
func (*loggerMock) Close()                        {}
