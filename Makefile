all:
	mkdir -p ./peerr
	mkdir -p ./build
	go build -o ./peerr/main main.go
	cp -R ./templates ./peerr
	cp ./config.example.json ./peerr
	tar -czf ./build/peerr.tar.gz peerr
deploy:
	scp ./build/peerr.tar.gz root@atl.l.x6c.us:~/
	scp ./build/peerr.tar.gz root@fra.l.x6c.us:~/
	scp ./build/peerr.tar.gz root@fnc.l.x6c.us:~/
	scp ./build/peerr.tar.gz root@lhr.l.x6c.us:~/
	scp ./build/peerr.tar.gz root@nrt.l.x6c.us:~/
	scp ./build/peerr.tar.gz root@ewr.l.x6c.us:~/
	scp ./build/peerr.tar.gz root@sea.v.x6c.us:~/

clean:
	rm -Rf build
	rm -Rf peerr
.PHONY: all clean
