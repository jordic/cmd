
del_stopped(){
    local name=$1
    local state=$(docker inspect --format "{{.State.Running}}" $name 2>/dev/null)

    if [[ "$state" == "false" ]]; then
        docker rm $name
    fi
}


mysql_server_install() {
    docker run --name mysql -v /home/jordi/mysql:/var/lib/mysql -p 3306:3306 -d mysql:5.6
}

mysql_start(){
    echo "Stargint mysql container"
    docker start mysql
}

mysql_stop() {
    docker stop mysql
}

redis_cli() {
    docker run -it --link redis:redis --rm redis sh -c 'exec redis-cli -h "$REDIS_PORT_6379_TCP_ADDR" -p "$REDIS_PORT_6379_TCP_PORT"'
}

redis_install() {
    docker run --name redis -p 6379:6379 -d redis
}
redis_start() {
    docker start redis
}

# gulp and node
gu() {
    docker run --rm -it -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp $@
}

liveserver() {

    docker run --rm -it -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) -p 8080:8080 jordic/gulp live-server
}


gulp() {
    docker run --rm -it -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp gulp $@
}

npm() {
  docker run --rm -it -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp npm $@
}

yo() {
  docker run -it --rm -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp yo $@
}

_jshint() {
  docker run --rm -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp jshint $@
}
bower() {
  docker run --rm -it -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp bower $@
}


build_chrome() {
    xhost +    
    docker run -it --net host -v /tmp/.X11-unix:/tmp/.X11-unix \
        -e DISPLAY=unix$DISPLAY -v $HOME/Downloads:/root/Downloads \
        -v $HOME/.config/google-chrome/:/data \
        --device /dev/snd -v /dev/shm:/dev/shm --name chrome jess/chrome


}



nv() {
    docker run -it --rm -v $(pwd)/:/mnt/ \
        -e UID=$(id -u) \
        -e GID=$(id -g) \
        --name neovim \
        jordic/neovim
}
