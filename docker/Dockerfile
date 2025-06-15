FROM golang:latest

# 使用tool追加
RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    xvfb \
    xorg-dev \
    libgl1-mesa-dev \
    libopenal1 \
    libopenal-dev \
    libvorbis0a \
    libvorbis-dev \
    libvorbisfile3


# ユーザー作成
ARG USERNAME
ARG GROUPNAME
ARG UID
ARG GID
RUN groupadd -g $GID $GROUPNAME && \
    useradd -m -s /bin/bash -u $UID -g $GID -G sudo $USERNAME

# sudo したいとき
# RUN apt-get update \
#     && apt-get install -y sudo \
#     && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
#     && chmod 0440 /etc/sudoers.d/$USERNAME


# 追加でパッケージ追加したい場合はrootユーザーで実施すれば。。。（docker exec -it -u 0 $(docker ps -q -f "name=aws-q-sandbox") bash）
USER $USERNAME

# AWS Qをインストール
RUN curl --proto '=https' --tlsv1.2 -sSf "https://desktop-release.q.us-east-1.amazonaws.com/latest/q-x86_64-linux.zip" -o "/tmp/q.zip" \
    && unzip /tmp/q.zip -d /tmp \
    && /tmp/q/install.sh --no-confirm \
    && rm -rf /tmp/q  \
    && rm -f /tmp/q.zip
