default: upd_local

upd_local:
	go get -u github.com/isaac-weisberg/go-jason@v0.0.0-unpublished

# GOPROXY=proxy.golang.org go list -m github.com/isaac-weisberg/go-jason-gen@v0.1.1