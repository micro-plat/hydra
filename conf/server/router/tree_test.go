package router

import (
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func TestNode_Match(t *testing.T) {
	tests := []struct {
		services  []string
		path      string
		matched   bool
		matchPath string
	}{
		{services: []string{"/:a/:b/:c", "/a/b/c"}, path: "/a", matchPath: "", matched: false},
		{services: []string{"/:a/:b/:c", "/a/b/c"}, path: "/a/b", matchPath: "", matched: false},
		{services: []string{"/:a/:b/:c", "/a/b/c"}, path: "/a/b/c", matchPath: "/:a/:b/:c", matched: true},
		{services: []string{"/:a/:b/:c", "/a/b/c"}, path: "/a/b/d", matchPath: "/:a/:b/:c", matched: true},
		{services: []string{"/:a/:b/:c", "/a/b/c"}, path: "/a/b/c/d", matchPath: "", matched: false},
		{services: []string{"/a/b/c", "/a/b/c"}, path: "/a/params/c", matchPath: "", matched: false},
		{services: []string{"/a/b/c", "/a/b/c"}, path: "/a/params/c", matchPath: "", matched: false},
		{services: []string{"/a/b/c", "/a/b/c"}, path: "/a/params/c", matchPath: "", matched: false},
		{services: []string{"/a/b/c", "/a/:b/c"}, path: "/a/params/c", matchPath: "/a/:b/c", matched: true},
		{services: []string{"/a/b/c", "/a/:b/c"}, path: "/a/b/c", matchPath: "/a/b/c", matched: true},
		{services: []string{"/a/b/c", "/a/:b/:c"}, path: "/a/param/param", matchPath: "/a/:b/:c", matched: true},
		{services: []string{"/a/b/c", "/a/:b/:c"}, path: "/a/param/param/params", matchPath: "", matched: false},
		{services: []string{"/a/b/c", "/:a/:b/d"}, path: "/param/param/d", matchPath: "/:a/:b/d", matched: true},
		{services: []string{"/a/b/c", "/:a/:b/c"}, path: "/param/param/d", matchPath: "", matched: false},
		{services: []string{"/a/b/c", "/a/b/e"}, path: "/a/b/d", matchPath: "", matched: false},
		{services: []string{"/a/b/c", "/:a/b/:c"}, path: "/param/b/param", matchPath: "/:a/b/:c", matched: true},
		{services: []string{"/a/b/c", "/a/:b/c", "/a/:b/d"}, path: "/a/b/d", matchPath: "/a/:b/d", matched: true},
		{services: []string{"/a/:b/c", "/a/:b/d"}, path: "/a/b/c", matchPath: "/a/:b/c", matched: true},
		{services: []string{"/a/:b/c", "/a/:b/d", "a/b/c"}, path: "/a/b/c", matchPath: "/a/:b/c", matched: true},
		{services: []string{"a/b/c", "/a/:b/c", "/a/:b/d"}, path: "/a/b/c", matchPath: "/a/b/c", matched: true},
		{services: []string{}, path: "/a/b/c", matchPath: "", matched: false},
		{services: []string{}, path: "/", matchPath: "", matched: false},
		{services: []string{"/", "/a"}, path: "/", matchPath: "/", matched: true},
	}
	for index, tt := range tests {
		tree := NewTree(tt.services...)
		matchPath, matched := tree.Match(tt.path, "")
		assert.Equal(t, tt.matched, matched, index)
		assert.Equal(t, tt.matchPath, matchPath, index)
	}
}
