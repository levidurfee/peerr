all:
	mkdir -p ./peerr
	mkdir -p ./build
	go build -o ./peerr/main main.go
	cp -R ./templates ./peerr
	cp ./config.example.json ./peerr
	tar -czf ./build/peerr.tar.gz peerr
deploy:
	scp ./build/peerr.tar.gz root@fra.x6c.us:~/
	scp ./build/peerr.tar.gz root@fnc.x6c.us:~/
	scp ./build/peerr.tar.gz root@lhr.x6c.us:~/
	scp ./build/peerr.tar.gz root@nrt.x6c.us:~/
	scp ./build/peerr.tar.gz root@ewr.x6c.us:~/
	scp ./build/peerr.tar.gz root@sgp.x6c.us:~/

clean:
	rm -Rf build
	rm -Rf peerr
.PHONY: all clean
