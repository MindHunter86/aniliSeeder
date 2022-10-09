# aniliSeeder
[![DeepSource](https://deepsource.io/gh/MindHunter86/aniliSeeder.svg/?label=active+issues&show_trend=true&token=0s6kHn6xfivpVWxqql7PLY23)](https://deepsource.io/gh/MindHunter86/aniliSeeder/?ref=repository-badge)


## Project is not ready now! Not for production use!!
---

<br/>

## cli internal commands (cli socket located on master)
`getTorrents` - get current torrents with some detailed info from all workers (cached data from conection init phase)

`listWorkers` - get all connected workers (and not connected #30 :))

`aniUpdates` - get last 5 anime updates from the anilibria api

`aniChanges` - get last 5 anime changes from the anilibria api

`aniSchedule` - get the current anilibria team shedule


## Running
run worker - `./aniliSeeder --http-debug --grpc-insecure serve`

run master - `./aniliSeeder --http-debug --grpc-insecure --swarm-is-master serve`

run cli = `./aniliSeeder cli`

## Usage
```
NAME:
   aniliSeeder - N\A

USAGE:
   aniliSeeder [global options] command [command options] [arguments...]

VERSION:
   v0.1

AUTHOR:
   MindHunter86 <admin@vkom.cc>

COMMANDS:
   serve    
   cli     
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --anilibria-api-baseurl value       (default: "https://api.anilibria.tv/v2")
   --anilibria-baseurl value           (default: "https://www.anilibria.tv")
   --anilibria-login-password value    password [$ANILIBRIA_PASSWORD]
   --anilibria-login-username value    login [$ANILIBRIA_LOGIN, $ANILIBRIA_USERNAME]
   --deluge-addr value                 (default: "127.0.0.1:58846")
   --deluge-data-path value            (default: "./data")
   --deluge-password value              [$DELUGE_PASSWORD]
   --deluge-torrentfiles-path value    (default: "./data")
   --deluge-username value             (default: "localclient") [$DELUGE_LOGIN, $DELUGE_USERNAME]
   --disk-minimal-avaliable value      In MB (default: 128)
   --grpc-connect-timeout value        for worker (default: 3s)
   --grpc-disable-reconnect            (default: false)
   --grpc-insecure                     (default: false)
   --grpc-ping-interval value          0 for disabling (default: 1s)
   --grpc-ping-reconnect-hold value    time for grpc reconnection process (default: 10s)
   --grpc-request-timeout value        (default: 1s)
   --help, -h                          show help (default: false)
   --http-client-insecure              Flag for TLS certificate verification disabling (default: false)
   --http-client-timeout TIMEOUT       Internal HTTP client connection TIMEOUT (format: 1000ms, 1s) (default: 3s)
   --http-debug                        (default: false)
   --http-idle-timeout value           (default: 5m0s)
   --http-keepalive-timeout value      (default: 5m0s)
   --http-max-idle-conns value         (default: 100)
   --http-tcp-timeout value            (default: 1s)
   --http-tls-handshake-timeout value  (default: 1s)
   --http2-conn-max-age value          for master; 0 for disable (default: 10m0s)
   --http2-ping-time value             for worker (default: 3s)
   --http2-ping-timeout value          for worker (default: 1s)
   --print-version, -V                 (default: false)
   --quite, -q                         Flag is equivalent to verbose -1 (default: false)
   --socket-path value                 (default: "aniliSeeder.sock")
   --swarm-custom-ca-path value        
   --swarm-is-master                   (default: false)
   --swarm-master-addr value           (default: "localhost:8081")
   --swarm-master-listen value         (default: "localhost:8081")
   --swarm-master-secret value         (default: "randomsecretkey") [$SWARM_MASTER_SECRETKEY]
   --torrentfiles-dir value            (default: "./data")
   --torrents-vkscore-line value       (default: 25)
   --verbose LEVEL, -v LEVEL           Verbose LEVEL (value from 5(debug) to 0(panic) and -1 for log disabling(quite mode)) (default: 5)
   

COPYRIGHT:
   (c) 2022 mindhunter86
```