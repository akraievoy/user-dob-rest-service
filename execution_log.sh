ak@pewznl:~/w/tsv_load$ ./00_bootstrap.sh 
you might pass 'force' as first param to this script to force dependency/protobuf re-fetch
STEP 1: installing git and protobuf (in golang docker)
fetch http://dl-cdn.alpinelinux.org/alpine/v3.9/main/x86_64/APKINDEX.tar.gz
fetch http://dl-cdn.alpinelinux.org/alpine/v3.9/community/x86_64/APKINDEX.tar.gz
v3.9.1-5-gae3d795467 [http://dl-cdn.alpinelinux.org/alpine/v3.9/main]
v3.9.0-90-gf0b6b9da2f [http://dl-cdn.alpinelinux.org/alpine/v3.9/community]
OK: 9754 distinct packages available
OK: 6 MiB in 15 packages
(1/9) Installing nghttp2-libs (1.35.1-r0)
(2/9) Installing libssh2 (1.8.0-r4)
(3/9) Installing libcurl (7.64.0-r1)
(4/9) Installing expat (2.2.6-r0)
(5/9) Installing pcre2 (10.32-r1)
(6/9) Installing git (2.20.1-r0)
(7/9) Installing libgcc (8.2.0-r2)
(8/9) Installing libstdc++ (8.2.0-r2)
(9/9) Installing protobuf (3.6.1-r1)
Executing busybox-1.29.3-r10.trigger
OK: 26 MiB in 24 packages
STEP 2: installing protoc-gen-go
STEP 3: generating protobuf golang API
STEP 4: vendoring library dependencies
go: finding github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51
go: finding github.com/golang/protobuf v1.3.0
go: finding github.com/facebookgo/flagenv v0.0.0-20160425205200-fcd59fca7456
go: finding github.com/davecgh/go-spew v1.1.1
go: finding github.com/lib/pq v1.0.0
go: finding github.com/facebookgo/stack v0.0.0-20160209184415-751773369052
go: finding github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870
go: finding google.golang.org/grpc v1.19.0
go: finding google.golang.org/genproto v0.0.0-20180831171423-11092d34479b
go: finding golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f
go: finding golang.org/x/net v0.0.0-20180906233101-161cd47e91fd
go: finding golang.org/x/net v0.0.0-20180826012351-8a410e7b638d
go: finding github.com/golang/protobuf v1.2.0
go: finding github.com/client9/misspell v0.3.4
go: finding github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
go: finding google.golang.org/appengine v1.1.0
go: finding golang.org/x/tools v0.0.0-20190114222345-bf090417da8b
go: finding golang.org/x/sys v0.0.0-20180830151530-49385e6e1522
go: finding golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
go: finding google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8
go: finding golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3
go: finding golang.org/x/text v0.3.0
go: finding cloud.google.com/go v0.26.0
go: finding honnef.co/go/tools v0.0.0-20190102054323-c2f93a96b099
go: finding github.com/BurntSushi/toml v0.3.1
go: finding github.com/golang/mock v1.1.1
go: downloading github.com/golang/protobuf v1.3.0
go: downloading github.com/facebookgo/flagenv v0.0.0-20160425205200-fcd59fca7456
go: downloading github.com/lib/pq v1.0.0
go: downloading google.golang.org/grpc v1.19.0
go: extracting github.com/facebookgo/flagenv v0.0.0-20160425205200-fcd59fca7456
go: extracting github.com/lib/pq v1.0.0
go: extracting github.com/golang/protobuf v1.3.0
go: extracting google.golang.org/grpc v1.19.0
go: downloading google.golang.org/genproto v0.0.0-20180831171423-11092d34479b
go: downloading golang.org/x/net v0.0.0-20180906233101-161cd47e91fd
go: downloading golang.org/x/sys v0.0.0-20180830151530-49385e6e1522
go: extracting golang.org/x/net v0.0.0-20180906233101-161cd47e91fd
go: extracting golang.org/x/sys v0.0.0-20180830151530-49385e6e1522
go: downloading golang.org/x/text v0.3.0
go: extracting google.golang.org/genproto v0.0.0-20180831171423-11092d34479b
go: extracting golang.org/x/text v0.3.0
docker may have generated some root-owned files -- reclaim ownership?
ak@pewznl:~/w/tsv_load$ ./01_build_dockers.sh
STEP 1: compiling ingestor
STEP 2: compiling upserter
STEP 3: purging previously built dockers
No stopped containers
STEP 4: purging previously built and now untagged docker images
Deleted: sha256:5403cff6737baae8b1017efd61db6982dfa892189b27d71a2e3ed7b5ad441588
Deleted: sha256:11d091007a10fde39c2469a5f964eae2289223c4804c689b90d0317a845b5b26
Deleted: sha256:c831e305ef44ca498816c4e4bdeddf7ad7dab6eecc1812bb968d40286b850c26
Deleted: sha256:70a137fe75413b2cba74c5c69df8908adfc13b127493c0cb099561c1a6c14123
Deleted: sha256:13a4663818093cc9b11b5db79155918ee876fde1a41077c99e971607ebef1c6b
Deleted: sha256:5c569ef758a4e6b5738a8266d400822de569a9ae4932b8a6b76bb14156e28651
Deleted: sha256:df477c05c29d59bae8bfd391af7566438d5c6b27a7e9c3c939d150dfde3802f6
Deleted: sha256:81f03b61cd91df39921ef98fe8b861f49c0cd4ec5388eeb3cafb51c8d28bd869
Deleted: sha256:d7b2731c164473bf2156828bca063ed0bc6831aa411cf17e96894d512cafda70
Deleted: sha256:c4213741e4e2e416ef763a36dfae56c030ef560c81c253aff1b1dab4c1cc61fb
Deleted: sha256:a1ecb28bb3612f7269949e162377232b7ad057123570f2838cb426cacf42f9ad
Deleted: sha256:66719d5bb8492c5faa60dc5b8a7b6f3ef1e76927daed4d1382294cb0a3ba6427
Deleted: sha256:122abc4e37ff04eadc9f793dd34fc5761736fce9a69fe757a86a8852bb187b34
Total reclaimed space: 0B
STEP 5: building dockers
verifier uses an image, skipping
postgres uses an image, skipping
Building upserter
Step 1/10 : FROM alpine:3.9
 ---> caf27325b298
Step 2/10 : COPY upserter /
 ---> Using cache
 ---> 093ad55ef7f0
Step 3/10 : ENV BIND_HOST=upserter
 ---> Using cache
 ---> dc32decd3ab2
Step 4/10 : ENV BIND_PORT=8080
 ---> Using cache
 ---> 9d37b2b7551f
Step 5/10 : ENV DB_HOST=postgres
 ---> Using cache
 ---> 108e3514432a
Step 6/10 : ENV DB_PORT=5432
 ---> Using cache
 ---> a87b60ddd9a5
Step 7/10 : ENV DB_USERNAME=upserter
 ---> Using cache
 ---> f0f14c99cdd8
Step 8/10 : ENV DB_PASSWORD=fillmein
 ---> Using cache
 ---> 83c31c82f7b2
Step 9/10 : ENV DB_NAME=upserter
 ---> Using cache
 ---> 2b65775fe069
Step 10/10 : CMD [   "/upserter" ]
 ---> Using cache
 ---> 48958a3e1d70
Successfully built 48958a3e1d70
Successfully tagged upserter:latest
Building ingestor
Step 1/9 : FROM alpine:3.9
 ---> caf27325b298
Step 2/9 : COPY ingestor /
 ---> 68540a095544
Step 3/9 : ENV IN_FILE_PATH=localhost
 ---> Running in 605979d51c9e
Removing intermediate container 605979d51c9e
 ---> f8178f7cdd35
Step 4/9 : ENV BROKEN_LINES_TO_FAIL=64
 ---> Running in 33bf938033aa
Removing intermediate container 33bf938033aa
 ---> 3b9e0c3829e9
Step 5/9 : ENV BATCH_SIZE=32
 ---> Running in 16fea094265d
Removing intermediate container 16fea094265d
 ---> c1b3149608a8
Step 6/9 : ENV UPSERTER_HOST=upserter
 ---> Running in bb3272583c3d
Removing intermediate container bb3272583c3d
 ---> ff8119f4fae8
Step 7/9 : ENV UPSERTER_PORT=8080
 ---> Running in 40059a2ab8ea
Removing intermediate container 40059a2ab8ea
 ---> 6f57ad40f9de
Step 8/9 : ENV FINAL_SLEEP=100ms
 ---> Running in 32d0cdf442be
Removing intermediate container 32d0cdf442be
 ---> bc466892c8cb
Step 9/9 : CMD [ "/ingestor" ]
 ---> Running in 0d4624700d75
Removing intermediate container 0d4624700d75
 ---> 23b76fafec64
Successfully built 23b76fafec64
Successfully tagged ingestor:latest
ak@pewznl:~/w/tsv_load$ ./02_run_end_to_end.sh 
STEP 1: creating/starting dockers
Creating network "tsv_load_default" with the default driver
Creating verifier ... done
Creating postgres ... done
Creating upserter ... done
Creating ingestor ... done
Attaching to verifier, postgres, upserter, ingestor
verifier    | psql: could not connect to server: Connection refused
verifier    | 	Is the server running on host "postgres" (172.18.0.2) and accepting
verifier    | 	TCP/IP connections on port 5432?
verifier    | psql: could not connect to server: Connection refused
verifier    | 	Is the server running on host "postgres" (172.18.0.2) and accepting
verifier    | 	TCP/IP connections on port 5432?
verifier    | sleeping on query 0: status=2, retries=1
ingestor    | 2019/03/05 09:55:10 dialling upserter:8080...
postgres    | The files belonging to this database system will be owned by user "postgres".
postgres    | This user must also own the server process.
postgres    | 
ingestor    | 2019/03/05 09:55:10 dialling upserter:8080 successful
ingestor    | 2019/03/05 09:55:10 failed to query upserter version (retrying in 5 sec): rpc error: code = Unavailable desc = all SubConns are in TransientFailure, latest connection error: connection error: desc = "transport: Error while dialing dial tcp 172.18.0.4:8080: connect: connection refused"
upserter    | 2019/03/05 09:55:10 database not available (will retry in 5 seconds): dial tcp 172.18.0.2:5432: connect: connection refused
postgres    | The database cluster will be initialized with locale "en_US.utf8".
postgres    | The default database encoding has accordingly been set to "UTF8".
postgres    | The default text search configuration will be set to "english".
postgres    | 
postgres    | Data page checksums are disabled.
postgres    | 
postgres    | fixing permissions on existing directory /var/lib/postgresql/data ... ok
postgres    | creating subdirectories ... ok
postgres    | selecting default max_connections ... 100
postgres    | selecting default shared_buffers ... 128MB
postgres    | selecting dynamic shared memory implementation ... posix
postgres    | creating configuration files ... ok
postgres    | running bootstrap script ... ok
postgres    | performing post-bootstrap initialization ... ok
postgres    | syncing data to disk ... 
postgres    | WARNING: enabling "trust" authentication for local connections
postgres    | You can change this by editing pg_hba.conf or using the option -A, or
postgres    | --auth-local and --auth-host, the next time you run initdb.
postgres    | ok
postgres    | 
postgres    | Success. You can now start the database server using:
postgres    | 
postgres    |     pg_ctl -D /var/lib/postgresql/data -l logfile start
postgres    | 
postgres    | waiting for server to start....LOG:  database system was shut down at 2019-03-05 09:55:11 UTC
postgres    | LOG:  MultiXact member wraparound protections are now enabled
postgres    | LOG:  database system is ready to accept connections
postgres    | LOG:  autovacuum launcher started
verifier    | psql: could not connect to server: Connection refused
verifier    | 	Is the server running on host "postgres" (172.18.0.2) and accepting
verifier    | 	TCP/IP connections on port 5432?
verifier    | sleeping on query 0: status=2, retries=2
postgres    |  done
postgres    | server started
postgres    | CREATE DATABASE
postgres    | 
postgres    | 
postgres    | /usr/local/bin/docker-entrypoint.sh: ignoring /docker-entrypoint-initdb.d/*
postgres    | 
postgres    | LOG:  received fast shutdown request
postgres    | waiting for server to shut down...LOG:  aborting any active transactions
postgres    | .LOG:  autovacuum launcher shutting down
postgres    | LOG:  shutting down
postgres    | LOG:  database system is shut down
postgres    |  done
postgres    | server stopped
postgres    | 
postgres    | PostgreSQL init process complete; ready for start up.
postgres    | 
postgres    | LOG:  database system was shut down at 2019-03-05 09:55:13 UTC
postgres    | LOG:  MultiXact member wraparound protections are now enabled
postgres    | LOG:  database system is ready to accept connections
postgres    | LOG:  autovacuum launcher started
upserter    | 2019/03/05 09:55:15 create table USER (if not exists), 0 rows affected
verifier    | sleeping on query 0: status=0, retries=3
ingestor    | 2019/03/05 09:55:15 failed to query upserter version (retrying in 5 sec): rpc error: code = Unavailable desc = all SubConns are in TransientFailure, latest connection error: connection error: desc = "transport: Error while dialing dial tcp 172.18.0.4:8080: connect: connection refused"
verifier    | query 0 PASSED
verifier    | sleeping on query 1: status=0, RES=     0, retries=4
ingestor    | 2019/03/05 09:55:20 connected to upserter version: 0.1.0
upserter    | 2019/03/05 09:55:21 upserted 32 records 
upserter    | 2019/03/05 09:55:21 upserted 32 records 
upserter    | 2019/03/05 09:55:21 upserted 32 records 
upserter    | 2019/03/05 09:55:21 upserted 7 records 
ingestor    | 2019/03/05 09:55:21 Completed: 103 records successfully processed, 0 records failed
ingestor    | 2019/03/05 09:55:21 sleeping for 1m30s
verifier    | sleeping on query 1: status=0, RES=   100, retries=5
verifier    | query 1 PASSED with RES=   100
verifier    | query 2 PASSED with RES=     1
verifier exited with code 0
Aborting on container exit...
Stopping ingestor ... done
Stopping upserter ... done
Stopping postgres ... done
END TO END PASSED
ak@pewznl:~/w/tsv_load$ mc

ak@pewznl:~/w/tsv_load$ ./03_clean.sh 
[sudo] password for ak: 
Sorry, try again.
[sudo] password for ak: 
chmod: cannot access '.idea': No such file or directory
Enumerating objects: 113, done.
Counting objects: 100% (113/113), done.
Delta compression using up to 2 threads
Compressing objects: 100% (64/64), done.
Writing objects: 100% (113/113), done.
Total 113 (delta 43), reused 113 (delta 43)
Going to remove ingestor, upserter, postgres, verifier
Removing ingestor ... done
Removing upserter ... done
Removing postgres ... done
Removing verifier ... done
Deleted: sha256:b51835afa4f1cf0ab27297afbafde0728b6cb0332b3d577911e8d20d2b9c3eb3
Deleted: sha256:76d2c9b2817596c6dc721b6c66dd93419d9f23492f158d2a112c1eb5e3888ff8
Deleted: sha256:47514cfb7284d22bd4b626ba0d8e821d1bc16cb7df687b8c1bd2413c9325648f
Deleted: sha256:ae7f12e219d7c7365a2eb5c50e6738c1aeeb894cf8a6e9df65a0be56e0af7144
Deleted: sha256:b77ed1ebbe79e9bc7244ebaa7ad76654c231151c2c197302d6cbf062bbf3af66
Deleted: sha256:f5a564e2959edf13fe92fdd274b334a0f73bab71b80d23e7a3aef00828b0cc85
Deleted: sha256:8aeb777a5f1cbe98830f0e5d18e86c2adaf9be85725948e81178c370eaaa076b
Deleted: sha256:c044a9484dff47c0e7d891eb4961797c0f57e872977abe6d6e6bc3d7eb48b597
Deleted: sha256:82650ae2ab688cbcda3949836f7f41dacd0c0b1d460de09100c611d851f12e54
Deleted Networks:
tsv_load_default

Deleted Volumes:
edd8519d9676054cdd68cf14b4417f70eb16c8f9f401125a605b5e9b16d77013
eb43e89c9d1dbcf63c3f2a8a09ad3b7f7933f6abecb03fe3c58bd0c29a693c04

Total reclaimed space: 46.09MB
ak@pewznl:~/w/tsv_load$ 

