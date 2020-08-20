FROM gocv-devel:v0.24.0-amd64 as devel
COPY * /root/go/src/github.com/baetyl/baetyl-video-infer/
RUN cd /root/go/src/github.com/baetyl/baetyl-video-infer/ && \
    make build

FROM debian:buster

RUN sed -i "s/deb.debian.org/mirrors.ustc.edu.cn/g" /etc/apt/sources.list && \
    sed -i "s/security.debian.org/mirrors.ustc.edu.cn/g" /etc/apt/sources.list && \
    apt-get update -y && \
    apt-get upgrade -y && \
    apt-get -y --no-install-recommends install pkg-config libgtk2.0-dev libavcodec-dev libavformat-dev libswscale-dev libtbb2 libtbb-dev libjpeg-dev libpng-dev libtiff-dev libdc1394-22-dev

# OpenCV
COPY --from=devel /usr/local/lib /usr/local/lib
COPY --from=devel /usr/local/lib/pkgconfig/opencv4.pc /usr/local/lib/pkgconfig/opencv4.pc
COPY --from=devel /usr/local/include/opencv4/opencv2 /usr/local/include/opencv4/opencv2

ENV PKG_CONFIG_PATH /usr/local/lib/pkgconfig
ENV LD_LIBRARY_PATH /usr/local/lib

# baetyl-video-infer
COPY --from=devel /root/go/src/github.com/baetyl/baetyl-video-infer/baetyl-video-infer /bin/

ENTRYPOINT ["baetyl-video-infer"]