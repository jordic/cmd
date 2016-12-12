



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

mysql_ram_install() {
  sudo cp -R /home/jordi/mysql /dev/shm/mysql
  docker run --name mysql -v /dev/shm/mysql:/var/lib/mysql -p 3306:3306 -d mysql:5.6
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

samba_create() {
    docker run -it --name samba -p 139:139 -p 445:445 -p 138:138/udp -p 445:445/udp \
        -e SMB_USER='jordi' -e SMB_PASS='tempo' -v /home/jordi:/home -d appcontainers/samba
}


# gulp and node
gu() {
    docker run --rm -it -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp $@
}

dliveserver() {

    docker run --rm -it -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) -p 8080:8080 jordic/gulp live-server
}


dgulp() {
    docker run --rm -it -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp gulp $@
}

dnpm() {
  docker run --rm -it -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp npm $@
}

dkarma() {
  docker run --rm -it -v $(pwd)/:/mnt/ -p 9876:9876 -e UID=$(id -u) -e GID=$(id -g) jordic/gulp karma $@
}

dyo() {
  docker run -it --rm -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp yo $@
}

_jshint() {
  docker run --rm -v $(pwd)/:/mnt/ -e UID=$(id -u) -e GID=$(id -g) jordic/gulp jshint $@
}
dbower() {
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
    del_stopped neovim
    docker run -it --rm -v $(pwd)/:/mnt/ \
        -e UID=$(id -u) \
        -e GID=$(id -g) \
        --name neovim \
        jordic/neovim
}


docker_net() {
    CID=$(docker ps -q --filter="name=${1}")
    echo -e "Container ID: $CID"
    docker inspect --format='{{ .NetworkSettings.Networks.tempo.IPAddress }}' $CID
}

sql_dump_docker() {
    CMD="mysqldump -h mysql -u tempo -p${DB_PASS} ${2}"
    ssh -C $1 "docker exec mysql ${CMD}"
}

sql_dump() {
    ssh -C $1 "mysqldump -u tempo -p${DB_PASS} ${2}"
}

tcookie() {
    cookiecutter https://github.com/tmpo-io/djangocms_base 
} 


wordpress() {
  
  if [ -z "$1" ]
  then
    name=wptest
  else
    name=$1
  fi

  docker rm $name
  docker run --name $name --net tempo \
    -e WORDPRESS_DB_HOST=mysql \
    -e WORDPRESS_DB_USER=wp2 \
    -e WORDPRESS_DB_PASSWORD=wp \
    -e WORDPRESS_DB_NAME=$name \
    -p 8080:80 -v $(pwd)/$name:/var/www/html \
    -d wordpress
}







