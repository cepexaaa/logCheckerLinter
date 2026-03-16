package zap

// stopper
type baseLogger struct{}

func (baseLogger) Info(...interface{})           {}
func (baseLogger) Infof(string, ...interface{})  {}
func (baseLogger) Infoln(...interface{})         {}
func (baseLogger) Error(...interface{})          {}
func (baseLogger) Errorf(string, ...interface{}) {}
func (baseLogger) Errorln(...interface{})        {}
func (baseLogger) Debug(...interface{})          {}
func (baseLogger) Debugf(string, ...interface{}) {}
func (baseLogger) Debugln(...interface{})        {}
func (baseLogger) Warn(...interface{})           {}
func (baseLogger) Warnf(string, ...interface{})  {}
func (baseLogger) Warnln(...interface{})         {}
func (baseLogger) Fatal(...interface{})          {}
func (baseLogger) Fatalf(string, ...interface{}) {}
func (baseLogger) Fatalln(...interface{})        {}
func (baseLogger) Panic(...interface{})          {}
func (baseLogger) Panicf(string, ...interface{}) {}
func (baseLogger) Panicln(...interface{})        {}

type Logger struct{ baseLogger }
type SugarLogger struct{ baseLogger }

func NewProduction(...interface{}) (*Logger, error) { return &Logger{}, nil }
func (l *Logger) Sugar() *SugarLogger               { return &SugarLogger{} }
func (s *SugarLogger) Desugar() *Logger             { return &Logger{} }
