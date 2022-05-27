package log

import "log"

func Panic(label, traceid string, data interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Panic(format)
		return
	}
	logger.Panic(format)
}

func Panicf(label, traceid string, data interface{}, args ...interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Panicf(format)
		return
	}
	logger.Panicf(format, args...)
}

func Fatal(label, traceid string, data interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Fatal(format)
		return
	}
	logger.Fatal(format)
}

func Fatalf(label, traceid string, data interface{}, args ...interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Fatalf(format)
		return
	}
	logger.Fatalf(format, args...)
}

func Debug(label, traceid string, data interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Print(format)
		return
	}
	logger.Debug(format)
}

func Debugf(label, traceid string, data interface{}, args ...interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Printf(format)
		return
	}
	logger.Debugf(format, args...)
}

func Info(label, traceid string, data interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Print(format)
		return
	}
	logger.Info(format)
}

func Infof(label, traceid string, data interface{}, args ...interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Printf(format)
		return
	}
	logger.Infof(format, args...)
}

func Error(label, traceid string, data interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Print(format)
		return
	}
	logger.Error(format)
}

func Errorf(label, traceid string, data interface{}, args ...interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Printf(format)
		return
	}
	logger.Errorf(format, args...)
}

func Warn(label, traceid string, data interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Print(format)
		return
	}
	logger.Warn(format)
}

func Warnf(label, traceid string, data interface{}, args ...interface{}) {
	format := format(label, traceid, data)
	if logger == nil {
		log.Printf(format)
		return
	}
	logger.Warnf(format, args...)
}
