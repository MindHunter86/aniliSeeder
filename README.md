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

`dryDeployAniUpdates` - check the titles (and their assignments) that can be deployed to workers

`deployAniUpdates` - deploy titles from `aniUpdates`

`dryDeployFailedAnnounces` - check worker for torrents with failed anounces to the anilibria tracker

`deployFailedAnnounces` - redeploy torrents with failed announces


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
   devel

AUTHOR:
   MindHunter86 <admin@vkom.cc>

COMMANDS:
   cli      
   serve    
   test     
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --anilibria-api-baseurl value       (default: "https://api.anilibria.tv/v2")
   --anilibria-baseurl value           (default: "https://www.anilibria.tv")
   --anilibria-login-password value    password [$ANILIBRIA_PASSWORD]
   --anilibria-login-username value    login [$ANILIBRIA_USERNAME]
   --cmd-vkscore-warn value            all torrents below this value will be marked as inefficient (default: 25)
   --cron-disable                      (default: false)
   --deluge-addr value                 (default: "127.0.0.1:58846")
   --deluge-data-path value            directory for space monitoring (default: "./data")
   --deluge-disk-minimal value         in MB;  (default: 128)
   --deluge-password value              [$DELUGE_PASSWORD]
   --deluge-torrents-path value        download directory for .torrent files (default: "./data")
   --deluge-username value             (default: "localclient") [$DELUGE_USERNAME]
   --deploy-ignore-errors              (default: false)
   --grpc-connect-timeout value        for worker (default: 3s)
   --grpc-disable-reconnect            (default: false)
   --grpc-insecure                     (default: false)
   --grpc-ping-interval value          0 for disabling (default: 1s)
   --grpc-ping-reconnect-hold value    time for grpc reconnection process (default: 5s)
   --grpc-reconnect-tries value        (default: 10)
   --grpc-request-timeout value        (default: 1s)
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
   --is-master                         (default: false) [$IS_MASTER]
   --master-addr value                 (default: "localhost:8081")
   --master-mon-interval value         master workers monitoring checks interval; 0 - for disabling (default: 3s)
   --master-secret value               (default: "randomsecretkey") [$SWARM_MASTER_SECRETKEY]
   --quite, -q                         Flag is equivalent to verbose -1 (default: false)
   --socket-path value                 (default: "aniliSeeder.sock")
   --syslog-addr value                 (default: "10.10.11.1:33517")
   --syslog-proto value                (default: "tcp")
   --syslog-tag value                  (default: "aniliseeder")
   --verbose LEVEL, -v LEVEL           Verbose LEVEL (value from 5(debug) to 0(panic) and -1 for log disabling(quite mode)) (default: 5)
   --help, -h                          show help (default: false)
   --print-version, -V                 (default: false)

COPYRIGHT:
   (c) 2022 mindhunter86
```
