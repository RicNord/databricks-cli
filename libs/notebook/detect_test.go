package notebook

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/databricks/databricks-sdk-go/service/workspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetect(t *testing.T) {
	var nb bool
	var lang workspace.Language
	var err error

	nb, lang, err = Detect("./testdata/py.py")
	require.NoError(t, err)
	assert.True(t, nb)
	assert.Equal(t, workspace.LanguagePython, lang)

	nb, lang, err = Detect("./testdata/r.r")
	require.NoError(t, err)
	assert.True(t, nb)
	assert.Equal(t, workspace.LanguageR, lang)

	nb, lang, err = Detect("./testdata/scala.scala")
	require.NoError(t, err)
	assert.True(t, nb)
	assert.Equal(t, workspace.LanguageScala, lang)

	nb, lang, err = Detect("./testdata/sql.sql")
	require.NoError(t, err)
	assert.True(t, nb)
	assert.Equal(t, workspace.LanguageSql, lang)

	nb, lang, err = Detect("./testdata/txt.txt")
	require.NoError(t, err)
	assert.False(t, nb)
	assert.Equal(t, workspace.Language(""), lang)
}

func TestDetectUnknownExtension(t *testing.T) {
	nb, _, err := Detect("./testdata/doesntexist.foobar")
	require.NoError(t, err)
	assert.False(t, nb)
}

func TestDetectNoExtension(t *testing.T) {
	nb, _, err := Detect("./testdata/doesntexist")
	require.NoError(t, err)
	assert.False(t, nb)
}

func TestDetectFileDoesNotExists(t *testing.T) {
	_, _, err := Detect("./testdata/doesntexist.py")
	require.Error(t, err)
}

func TestDetectEmptyFile(t *testing.T) {
	// Create empty file.
	dir := t.TempDir()
	path := filepath.Join(dir, "file.py")
	err := os.WriteFile(path, nil, 0644)
	require.NoError(t, err)

	// No contents means not a notebook.
	nb, _, err := Detect(path)
	require.NoError(t, err)
	assert.False(t, nb)
}

func TestDetectFileWithLongHeader(t *testing.T) {
	// Create 128kb garbage file.
	dir := t.TempDir()
	path := filepath.Join(dir, "file.py")
	buf := make([]byte, 128*1024)
	err := os.WriteFile(path, buf, 0644)
	require.NoError(t, err)

	// Garbage contents means not a notebook.
	nb, _, err := Detect(path)
	require.NoError(t, err)
	assert.False(t, nb)
}