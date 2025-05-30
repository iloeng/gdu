package common_test

import (
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/dundee/gdu/v5/internal/common"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

func TestCreateIgnorePattern(t *testing.T) {
	re, err := common.CreateIgnorePattern([]string{"[abc]+"})

	assert.Nil(t, err)
	assert.True(t, re.MatchString("aa"))
}

func TestCreateIgnorePatternWithErr(t *testing.T) {
	re, err := common.CreateIgnorePattern([]string{"[[["})

	assert.NotNil(t, err)
	assert.Nil(t, re)
}

func TestEmptyIgnore(t *testing.T) {
	ui := &common.UI{}
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.False(t, shouldBeIgnored("abc", "/abc"))
	assert.False(t, shouldBeIgnored("xxx", "/xxx"))
}

func TestIgnoreByAbsPath(t *testing.T) {
	ui := &common.UI{}
	ui.SetIgnoreDirPaths([]string{"/abc"})
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("abc", "/abc"))
	assert.False(t, shouldBeIgnored("xxx", "/xxx"))
}

func TestIgnoreByPattern(t *testing.T) {
	ui := &common.UI{}
	err := ui.SetIgnoreDirPatterns([]string{"/[abc]+"})
	assert.Nil(t, err)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("aaa", "/aaa"))
	assert.True(t, shouldBeIgnored("aaa", "/aaabc"))
	assert.False(t, shouldBeIgnored("xxx", "/xxx"))
}

func TestIgnoreFromFile(t *testing.T) {
	file, err := os.OpenFile("ignore", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString("/aaa\n"); err != nil {
		panic(err)
	}
	if _, err := file.WriteString("/aaabc\n"); err != nil {
		panic(err)
	}
	if _, err := file.WriteString("/[abd]+\n"); err != nil {
		panic(err)
	}

	ui := &common.UI{}
	err = ui.SetIgnoreFromFile("ignore")
	assert.Nil(t, err)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("aaa", "/aaa"))
	assert.True(t, shouldBeIgnored("aaabc", "/aaabc"))
	assert.True(t, shouldBeIgnored("aaabd", "/aaabd"))
	assert.False(t, shouldBeIgnored("xxx", "/xxx"))
}

func TestIgnoreFromNotExistingFile(t *testing.T) {
	ui := &common.UI{}
	err := ui.SetIgnoreFromFile("xxx")
	assert.NotNil(t, err)
}

func TestIgnoreHidden(t *testing.T) {
	ui := &common.UI{}
	ui.SetIgnoreHidden(true)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored(".git", "/aaa/.git"))
	assert.True(t, shouldBeIgnored(".bbb", "/aaa/.bbb"))
	assert.False(t, shouldBeIgnored("xxx", "/xxx"))
}

func TestIgnoreByAbsPathAndHidden(t *testing.T) {
	ui := &common.UI{}
	ui.SetIgnoreDirPaths([]string{"/abc"})
	ui.SetIgnoreHidden(true)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("abc", "/abc"))
	assert.True(t, shouldBeIgnored(".git", "/aaa/.git"))
	assert.True(t, shouldBeIgnored(".bbb", "/aaa/.bbb"))
	assert.False(t, shouldBeIgnored("xxx", "/xxx"))
}

func TestIgnoreByAbsPathAndPattern(t *testing.T) {
	ui := &common.UI{}
	ui.SetIgnoreDirPaths([]string{"/abc"})
	err := ui.SetIgnoreDirPatterns([]string{"/[abc]+"})
	assert.Nil(t, err)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("abc", "/abc"))
	assert.True(t, shouldBeIgnored("aabc", "/aabc"))
	assert.True(t, shouldBeIgnored("ccc", "/ccc"))
	assert.False(t, shouldBeIgnored("xxx", "/xxx"))
}

func TestIgnoreByPatternAndHidden(t *testing.T) {
	ui := &common.UI{}
	err := ui.SetIgnoreDirPatterns([]string{"/[abc]+"})
	assert.Nil(t, err)
	ui.SetIgnoreHidden(true)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("abbc", "/abbc"))
	assert.True(t, shouldBeIgnored(".git", "/aaa/.git"))
	assert.True(t, shouldBeIgnored(".bbb", "/aaa/.bbb"))
	assert.False(t, shouldBeIgnored("xxx", "/xxx"))
}

func TestIgnoreByAll(t *testing.T) {
	ui := &common.UI{}
	ui.SetIgnoreDirPaths([]string{"/abc"})
	err := ui.SetIgnoreDirPatterns([]string{"/[abc]+"})
	assert.Nil(t, err)
	ui.SetIgnoreHidden(true)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("abc", "/abc"))
	assert.True(t, shouldBeIgnored("aabc", "/aabc"))
	assert.True(t, shouldBeIgnored(".git", "/aaa/.git"))
	assert.True(t, shouldBeIgnored(".bbb", "/aaa/.bbb"))
	assert.False(t, shouldBeIgnored("xxx", "/xxx"))
}

func TestIgnoreByRelativePath(t *testing.T) {
	ui := &common.UI{}
	ui.SetIgnoreDirPaths([]string{"test_dir/abc"})
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("abc", "test_dir/abc"))
	absPath, err := filepath.Abs("test_dir/abc")
	assert.Nil(t, err)
	assert.True(t, shouldBeIgnored("abc", absPath))
	assert.False(t, shouldBeIgnored("xxx", "test_dir/xxx"))
}

func TestIgnoreByRelativePattern(t *testing.T) {
	ui := &common.UI{}
	err := ui.SetIgnoreDirPatterns([]string{"test_dir/[abc]+"})
	assert.Nil(t, err)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("abc", "test_dir/abc"))
	absPath, err := filepath.Abs("test_dir/abc")
	assert.Nil(t, err)
	assert.True(t, shouldBeIgnored("abc", absPath))
	assert.False(t, shouldBeIgnored("xxx", "test_dir/xxx"))
}

func TestIgnoreFromFileWithRelativePaths(t *testing.T) {
	file, err := os.OpenFile("ignore", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	defer os.Remove("ignore")

	if _, err := file.WriteString("test_dir/aaa\n"); err != nil {
		panic(err)
	}
	if _, err := file.WriteString("node_modules/[^/]+\n"); err != nil {
		panic(err)
	}

	ui := &common.UI{}
	err = ui.SetIgnoreFromFile("ignore")
	assert.Nil(t, err)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("aaa", "test_dir/aaa"))
	absPath, err := filepath.Abs("test_dir/aaa")
	assert.Nil(t, err)
	assert.True(t, shouldBeIgnored("aaa", absPath))
	assert.False(t, shouldBeIgnored("xxx", "test_dir/xxx"))
}
