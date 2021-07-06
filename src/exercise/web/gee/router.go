package gee

import (
    "fmt"
    "net/http"
    "reflect"
    "strings"
    "testing"
)

// Router implemented in Trie, which support high efficient dynamic routing
type Router struct {
    // roots: request method -> trie root
    roots map[string]*node
    // handlers: request method -> handler function
    handlers map[string]HandlerFunc
}

// newRouter: the Constructor of Router
func newRouter() *Router {
    return &Router{
        roots:    make(map[string]*node),
        handlers: make(map[string]HandlerFunc),
    }
}

// parsePattern: Only one * is allowed
func parsePattern(pattern string) []string {
    parts := make([]string, 0)
    for _, item := range strings.Split(pattern, "/") {
        if item != "" {
            parts = append(parts, item)

            // ignore the following parts
            if item[0] == '*' {
                break
            }
        }
    }
    return parts
}

// addRoute: add route to route map
func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
    parts := parsePattern(pattern)

    // GET-/p/:lang/doc
    key := method + "-" + pattern
    _, ok := r.roots[method]
    if !ok {
        r.roots[method] = &node{}
    }

    // GET -> trie, Get-/p/:lang/doc -> handler
    r.roots[method].insert(pattern, parts, 0)
    r.handlers[key] = handler
}

// getRoute: returns a trie node and a map stores params -> values
// for example: /static/css/geektutu.css matches /static/*filepath => {"filepath" -> "css/geektutu.css"}
func (r *Router) getRoute(method string, path string) (*node, map[string]string) {

    // /static/css/geektutu => ['static', 'css', 'geektutu']
    searchParts := parsePattern(path)
    params := make(map[string]string)
    root, ok := r.roots[method]
    if !ok {
        return nil, nil
    }

    n := root.search(searchParts, 0)
    if n != nil {
        parts := parsePattern(n.pattern)
        for index, part := range parts {
            if part[0] == ':' {
                params[part[1:]] = searchParts[index]
            }
            if part[0] == '*' && len(part) > 1 {
                params[part[1:]] = strings.Join(searchParts[index:], "/")
                break
            }
        }
        return n, params
    }

    return nil, nil
}

// getRoutes
func (r *Router) getRoutes(method string) []*node {
    root, ok := r.roots[method]
    if !ok {
        return nil
    }
    nodes := make([]*node, 0)
    root.travel(&nodes)
    return nodes
}

// handle handle request by get route from trie
func (r *Router) handle(c *Context) {
    n, params := r.getRoute(c.Method, c.Path)

    if n != nil {
        key := c.Method + "-" + n.pattern
        c.Params = params
        c.handlers = append(c.handlers, r.handlers[key])
    } else {
        c.handlers = append(c.handlers, func(c *Context) {
            c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
        })
    }
    c.Next()
}

func newTestRouter() *Router {
    r := newRouter()
    r.addRoute("GET", "/", nil)
    r.addRoute("GET", "/hello/:name", nil)
    r.addRoute("GET", "/hello/b/c", nil)
    r.addRoute("GET", "/hi/:name", nil)
    r.addRoute("GET", "/assets/*filepath", nil)
    return r
}

func TestParsePattern(t *testing.T) {
    ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
    ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
    ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
    if !ok {
        t.Fatal("test parsePattern failed")
    }
}

func TestGetRoute(t *testing.T) {
    r := newTestRouter()
    n, ps := r.getRoute("GET", "/hello/geektutu")

    if n == nil {
        t.Fatal("nil shouldn't be returned")
    }

    if n.pattern != "/hello/:name" {
        t.Fatal("should match /hello/:name")
    }

    if ps["name"] != "geektutu" {
        t.Fatal("name should be equal to 'geektutu'")
    }

    fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

}

func TestGetRoute2(t *testing.T) {
    r := newTestRouter()
    n1, ps1 := r.getRoute("GET", "/assets/file1.txt")
    ok1 := n1.pattern == "/assets/*filepath" && ps1["filepath"] == "file1.txt"
    if !ok1 {
        t.Fatal("pattern shoule be /assets/*filepath & filepath shoule be file1.txt")
    }

    n2, ps2 := r.getRoute("GET", "/assets/css/test.css")
    ok2 := n2.pattern == "/assets/*filepath" && ps2["filepath"] == "css/test.css"
    if !ok2 {
        t.Fatal("pattern shoule be /assets/*filepath & filepath shoule be css/test.css")
    }

}

func TestGetRoutes(t *testing.T) {
    r := newTestRouter()
    nodes := r.getRoutes("GET")
    for i, n := range nodes {
        fmt.Println(i+1, n)
    }

    if len(nodes) != 5 {
        t.Fatal("the number of routes shoule be 4")
    }
}
