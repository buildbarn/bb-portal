docker run -d \
  --name=sqlitebrowser \
  --security-opt seccomp=unconfined `#optional` \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=Etc/UTC \
  -p 4000:3000 \
  -p 4001:3001 \
  -v /path/to/config:/config \
  -v /data:/data \
  --restart unless-stopped \
  lscr.io/linuxserver/sqlitebrowser:latest