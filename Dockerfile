FROM centos:8

MAINTAINER actors315 <actors315@gmail.com>

ENV MYSQL_HOST=localhost
ENV MYSQL_PASSWORD=58117aec3b3252a97be0

# 基础环境
RUN mkdir -p /usr/local/src/qcloud-tools \
    && mkdir -p /usr/local/qcloud-tools/shell \
    && mkdir -p /usr/local/qcloud-tools/config \
    && yum install -y wget openssl make \
    && cd /usr/local/src \
	&& wget https://golang.google.cn/dl/go1.16.4.linux-amd64.tar.gz \
	&& wget -O acme.sh.tar.gz https://github.com/acmesh-official/acme.sh/archive/master.tar.gz \
# 安装 GO
	&& cd /usr/local/src \
	&& tar -C /usr/local -xvf go1.16.4.linux-amd64.tar.gz \
	&& echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bashrc \
	&& echo 'export GOROOT=/usr/local/go' >> ~/.bashrc \
	&& echo 'export GOPATH=/data/go' >> ~/.bashrc \
	&& echo 'export GOPROXY=https://mirrors.cloud.tencent.com/go/' >> ~/.bashrc \
# 安装 acme.sh
	&& cd /usr/local/src \
	&& tar -C /usr/local/src -xvf acme.sh.tar.gz \
	&& cd acme.sh-master \
	&& ./acme.sh --install --nocron

COPY . /usr/local/src/qcloud-tools/

# 安装 qcloud-tools
RUN cd /usr/local/src/qcloud-tools \
	&& make clean && make cert-monitor \
	&& mv /usr/local/src/qcloud-tools/bin/* /usr/local/qcloud-tools/ \
# 配置
	&& mv /usr/local/src/qcloud-tools/config/config.simple.yaml /usr/local/qcloud-tools/config/config.yaml \
    && mv /usr/local/src/qcloud-tools/config/issue-template.tpl /usr/local/qcloud-tools/config/issue-template.tpl \
    && mv /usr/local/src/qcloud-tools/web /usr/local/qcloud-tools/ \
    && mv /usr/local/src/qcloud-tools/Dockerstart /start \
    && chmod +x /start \
# 清理
    && rm -rf /usr/local/src/* \
    && yum remove -y wget \
    && yum clean all

EXPOSE 80

WORKDIR /usr/local/qcloud-tools/

CMD ["/start"]