package stash

import (
	"fmt"
	"github.com/adminhmi/hlog"
	"io"
	"sync"
)

// Hook represents a Logstash hook.
// It has two fields: writer to write the entry to Logstash and
// formatter to format the entry to a Logstash format before sending.
//
// To initialize it use the `New` function.
//
type Hook struct {
	writer    io.Writer
	formatter hlog.Formatter
}

// New returns a new logrus.Hook for Logstash.
//
// To create a new hook that sends logs to `tcp://logstash.corp.io:9999`:
//
// conn, _ := net.Dial("tcp", "logstash.corp.io:9999")
// hook := logrustash.New(conn, logrustash.DefaultFormatter())
func New(w io.Writer, f hlog.Formatter) hlog.Hook {
	return Hook{
		writer:    w,
		formatter: f,
	}
}

// Fire takes, formats and sends the entry to Logstash.
// Hook's formatter is used to format the entry into Logstash format
// and Hook's writer is used to write the formatted entry to the Logstash instance.
func (h Hook) Fire(e *hlog.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(dataBytes)
	return err
}

// Levels returns all hlog levels.
func (h Hook) Levels() []hlog.Level {
	return hlog.AllLevels
}

// Using a pool to re-use of old entries when formatting Logstash messages.
// It is used in the Fire function.
var entryPool = sync.Pool{
	New: func() interface{} {
		return &hlog.Entry{}
	},
}

func copyEntry(e *hlog.Entry, fields hlog.Fields) *hlog.Entry {
	ne := entryPool.Get().(*hlog.Entry)
	ne.Message = e.Message
	ne.Level = e.Level
	ne.Time = e.Time
	ne.Data = hlog.Fields{}

	if e.HasCaller() {
		ne.Data["function"] = e.Caller.Function
		ne.Data["file"] = fmt.Sprintf("%s:%d", e.Caller.File, e.Caller.Line)
	}

	for k, v := range fields {
		ne.Data[k] = v
	}
	for k, v := range e.Data {
		ne.Data[k] = v
	}
	return ne
}

// releaseEntry puts the given entry back to `entryPool`. It must be called if copyEntry is called.
func releaseEntry(e *hlog.Entry) {
	entryPool.Put(e)
}

type StashFormatter struct {
	hlog.Formatter
	hlog.Fields
}

var (
	stashFields   = hlog.Fields{"@version": "1", "type": "log"}
	stashFieldMap = hlog.FieldMap{
		hlog.FieldKeyTime: "@timestamp",
		hlog.FieldKeyMsg:  "message",
	}
)

// DefaultFormatter returns a default Logstash formatter:
// A JSON format with "@version" set to "1" (unless set differently in `fields`,
// "type" to "log" (unless set differently in `fields`),
// "@timestamp" to the log time and "message" to the log message.
//
// Note: to set a different configuration use the `LogstashFormatter` structure.
func DefaultFormatter(fields hlog.Fields) hlog.Formatter {
	for k, v := range stashFields {
		if _, ok := fields[k]; !ok {
			fields[k] = v
		}
	}

	return StashFormatter{
		Formatter: &hlog.JSONFormatter{FieldMap: stashFieldMap},
		Fields:    fields,
	}
}

// Format formats an entry to a Logstash format according to the given Formatter and Fields.
//
// Note: the given entry is copied and not changed during the formatting process.
func (f StashFormatter) Format(e *hlog.Entry) ([]byte, error) {
	ne := copyEntry(e, f.Fields)
	dataBytes, err := f.Formatter.Format(ne)
	releaseEntry(ne)
	return dataBytes, err
}
