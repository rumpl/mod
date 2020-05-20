package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MainSuite struct {
	tmpDir string
	suite.Suite
}

func (s *MainSuite) BeforeTest(suite, test string) {
	logrus.SetLevel(logrus.DebugLevel)
	dir, err := ioutil.TempDir("", "test")
	require.Nil(s.T(), err)

	s.tmpDir = dir

	err = os.Chdir(s.tmpDir)
	require.Nil(s.T(), err)
}

func (s *MainSuite) AfterTest(suite, test string) {
	err := os.RemoveAll(s.tmpDir)
	require.Nil(s.T(), err)
}

func (s *MainSuite) TestInit() {
	dir := "repo"

	err := run(dir)
	require.Nil(s.T(), err)

	if _, err := os.Stat(path.Join(s.tmpDir, dir)); os.IsNotExist(err) {
		s.Fail("Did not create repo dir")
	}

	if _, err := os.Stat(path.Join(s.tmpDir, dir, "README.md")); os.IsNotExist(err) {
		s.Fail("Did not create the README file")
	}
}

func (s *MainSuite) TestInitDirectoryExists() {
	dir := "repo"

	err := os.Mkdir(path.Join(s.tmpDir, dir), 0755)
	require.Nil(s.T(), err)

	err = run(dir)
	assert.Errorf(s.T(), err, "directory already exists")
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(MainSuite))
}
