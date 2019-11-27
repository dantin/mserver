package logutil

import (
	"bytes"
	"testing"

	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

const (
	logPattern = `\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}.\d{3} \[(fatal|error|warning|info|debug)\] .*?\n`
)

func Test(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testLogSuite{})

type testLogSuite struct {
	buf *bytes.Buffer
}

func (s *testLogSuite) SetUpSuite(c *C) {
	s.buf = &bytes.Buffer{}
}

func (s *testLogSuite) TestStringToLogLevel(c *C) {
	c.Assert(stringToLogLevel("fatal"), Equals, log.FatalLevel)
	c.Assert(stringToLogLevel("ERROR"), Equals, log.ErrorLevel)
	c.Assert(stringToLogLevel("warn"), Equals, log.WarnLevel)
	c.Assert(stringToLogLevel("warning"), Equals, log.WarnLevel)
	c.Assert(stringToLogLevel("debug"), Equals, log.DebugLevel)
	c.Assert(stringToLogLevel("info"), Equals, log.InfoLevel)
	c.Assert(stringToLogLevel("bad"), Equals, log.InfoLevel)
}

// TestLogging assure log format works.
func (s *testLogSuite) TestLogging(c *C) {
	conf := &LogConfig{Level: "warn", File: FileLogConfig{}}
	c.Assert(InitLogger(conf), IsNil)

	log.SetOutput(s.buf)

	log.Warnf("message from logrus")
	entry, err := s.buf.ReadString('\n')
	c.Assert(err, IsNil)
	c.Assert(entry, Matches, logPattern)
}
