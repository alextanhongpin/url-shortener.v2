package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var __000001_create_table_url_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x28\x2d\xca\xb1\xe6\x02\x04\x00\x00\xff\xff\x3c\x59\x05\x7d\x10\x00\x00\x00")

func _000001_create_table_url_down_sql() ([]byte, error) {
	return bindata_read(
		__000001_create_table_url_down_sql,
		"000001_create_table_url.down.sql",
	)
}

var __000001_create_table_url_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x8b\x4d\xaa\x83\x40\x10\x06\xd7\xd3\xa7\xf8\x96\x0a\xde\xe0\xad\xe6\x25\x2d\x69\xe2\x5f\x9c\x1e\xa2\xd9\x88\x64\x24\x04\x44\x83\x28\xb9\x7e\x88\x37\xc8\xb2\xa8\xaa\x43\xcd\x56\x19\x6a\xff\x33\x86\xa4\x28\x4a\x05\x37\xe2\xd4\x61\x5b\x46\x44\x64\x9e\x01\xca\x8d\xa2\xaa\x25\xb7\x75\x8b\x33\xb7\x09\x99\x71\x9e\x1e\xdd\xb7\xd8\x9d\x2f\xe4\xe2\x39\x21\x73\x5f\x86\x7e\x1d\x42\xd7\xaf\x50\xc9\xd9\xa9\xcd\x2b\x5c\x45\x4f\x3b\xe2\x56\x16\x8c\x23\xa7\xd6\x67\x8a\x69\x7e\x47\x71\x42\x66\x7b\x85\xdf\x1e\x8a\xff\xe8\x13\x00\x00\xff\xff\xeb\x5e\xc4\xf8\xb8\x00\x00\x00")

func _000001_create_table_url_up_sql() ([]byte, error) {
	return bindata_read(
		__000001_create_table_url_up_sql,
		"000001_create_table_url.up.sql",
	)
}

var _bindata_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func bindata_go() ([]byte, error) {
	return bindata_read(
		_bindata_go,
		"bindata.go",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"000001_create_table_url.down.sql": _000001_create_table_url_down_sql,
	"000001_create_table_url.up.sql": _000001_create_table_url_up_sql,
	"bindata.go": bindata_go,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"000001_create_table_url.down.sql": &_bintree_t{_000001_create_table_url_down_sql, map[string]*_bintree_t{
	}},
	"000001_create_table_url.up.sql": &_bintree_t{_000001_create_table_url_up_sql, map[string]*_bintree_t{
	}},
	"bindata.go": &_bintree_t{bindata_go, map[string]*_bintree_t{
	}},
}}
