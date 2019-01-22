
GLIDEBIN=${GOPATH}/bin/glide
GOFILES=$(find . -path ./vendor -prune -o -name '*.go' -print)

$(GLIDEBIN):
	go get github.com/Masterminds/glide

install: $(GLIDEBIN)
	${GOPATH}/bin/glide install

test: install
	go build
	go test -v ./
	# gofmt -l ${GOFILES} | read && echo "gofmt failures" && gofmt -d ${GOFILES} && exit 1 || true

statsd:
	brew list nodejs &>/dev/null || brew install nodejs
	git clone https://github.com/etsy/statsd.git
	echo '{ "port": 8125, "mgmt_port": 8126, "backends": [ "./backends/console" ] }' > statsd/config.json
	node statsd/stats.js statsd/config.json > /dev/null 2>&1 &

notifications:
	email:
		on_success: change
		on_failure: always


clean:
	rm -rf statsd