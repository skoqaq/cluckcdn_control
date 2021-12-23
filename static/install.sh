#!/bin/bash

# -----------
# - Install -
# -----------
source /etc/os-release
case $ID in
debian|ubuntu|devuan)
    sysinfo="apt"
    sudo apt update
    sudo apt install -y curl wget openssl build-essential libpcre3 libpcre3-dev zlib1g-dev git unzip
    ;;
centos|fedora|rhel)
    sysinfo="yum"
    if test "$(echo "$VERSION_ID >= 22" | bc)" -ne 0; then
        sysinfo="dnf"
    fi
    sudo $sysinfo -y install curl wget gcc pcre-devel zlib-devel git openssl openssl-devel unzip
    ;;
*)
    exit 1
    ;;
esac

# --------
# - Init -
# --------
sudo mkdir /node && sudo mkdir /node/error && sudo mkdir /node/cache && sudo mkdir /node/tls
sudo chmod 775 /node
# Download Node
wget -O /node/node '{ctrlServer}/static/node'
# Download Config
wget -O /node/config.yaml '{ctrlServer}/config.yaml?token={token}'
# Downloads Error Page
wget -O /node/error/error.zip '{ctrlServer}/static/error.zip'
unzip -d /node/error/ /node/error/error.zip
unzip -d /usr/local/nginx/html/ /node/error/error.zip

# -----------------
# - Install Nginx -
# -----------------

# Download Nginx
mkdir ~/nginx_build && cd ~/nginx_build
wget http://nginx.org/download/nginx-1.21.0.tar.gz
tar -xvf nginx-1.21.0.tar.gz
rm -rf nginx-1.21.0.tar.gz
# Download Br
git clone https://github.com/google/ngx_brotli.git
pushd ngx_brotli
git submodule update --init
popd
# Download openssl
wget -c https://www.openssl.org/source/openssl-1.1.1k.tar.gz
tar zxf openssl-1.1.1k.tar.gz
rm openssl-1.1.1k.tar.gz
# Cd Nginx
cd nginx-1.21.0
sed -i 's/Server: nginx/Server: CluckCDN/g' src/http/ngx_http_header_filter_module.c
# Setting
sudo ./configure \
--prefix=/usr/local/nginx \
--with-http_v2_module \
--with-http_ssl_module \
--with-http_gzip_static_module \
--with-http_sub_module \
--with-openssl=../openssl-1.1.1k \
--add-module=../ngx_brotli
# Make
sudo make
sudo make install
# Config
rm -f /usr/local/nginx/conf/nginx.conf
wget -O /usr/local/nginx/conf/nginx.conf '{ctrlServer}/static/nginx.conf'
wget -O /node/vhost.conf '{ctrlServer}/vhost.conf?token={token}'
# Start
/usr/local/nginx/sbin/nginx
/usr/local/nginx/sbin/nginx -s reload
# Start Node
cd /node && chmod 775 node && nohup ./node &