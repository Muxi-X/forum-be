ARG GO_IMAGE=golang:1.26.4
ARG PYTHON_IMAGE=python:3.10-slim-bookworm

FROM ${GO_IMAGE} AS builder

ARG GO_PROXY=https://goproxy.cn,direct

WORKDIR /app
COPY . /app/
RUN go env -w GOPROXY="${GO_PROXY}"

WORKDIR /app/microservice/ocr
RUN go mod tidy && go build -o /out/main .

FROM ${PYTHON_IMAGE}

ARG APT_MIRROR=
ARG APT_SECURITY_MIRROR=
ARG PIP_INDEX_URL=https://pypi.org/simple
ARG PYTORCH_INDEX_URL=https://download.pytorch.org/whl/cpu
ARG MODELSCOPE_FIND_LINKS=https://modelscope.oss-cn-beijing.aliyuncs.com/releases/repo.html

ENV TZ=Asia/Shanghai \
    HOME=/opt/forum-ocr \
    FORUM_OCR_HOME=/opt/forum-ocr \
    FORUM_OCR_VENV=/opt/forum-ocr/.venvs/modelscope-ocr \
    MODELSCOPE_CACHE=/opt/forum-ocr/modelscope-cache \
    PIP_ROOT_USER_ACTION=ignore \
    PYTHONUNBUFFERED=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

ARG PRELOAD_MODEL=true

RUN set -eux; \
    if [ -n "$APT_MIRROR" ]; then \
        find /etc/apt -type f \( -name '*.list' -o -name '*.sources' \) -print0 | \
        xargs -0 sed -i "s|http://deb.debian.org/debian|$APT_MIRROR|g; s|https://deb.debian.org/debian|$APT_MIRROR|g"; \
    fi; \
    if [ -n "$APT_SECURITY_MIRROR" ]; then \
        find /etc/apt -type f \( -name '*.list' -o -name '*.sources' \) -print0 | \
        xargs -0 sed -i "s|http://security.debian.org/debian-security|$APT_SECURITY_MIRROR|g; s|https://security.debian.org/debian-security|$APT_SECURITY_MIRROR|g"; \
    fi; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        libc6 \
        libglib2.0-0 \
        libsm6 \
        libxext6 \
        libxrender1 \
        tzdata; \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p \
    "$FORUM_OCR_HOME" \
    "$FORUM_OCR_HOME/.venvs" \
    "$FORUM_OCR_HOME/workspace" \
    "$MODELSCOPE_CACHE"

COPY microservice/ocr/requirements.txt /tmp/ocr-requirements.txt
COPY microservice/ocr/preload_model.py /tmp/preload_model.py

RUN set -eux; \
    python -m venv "$FORUM_OCR_VENV"; \
    "$FORUM_OCR_VENV/bin/pip" install --no-cache-dir --index-url "$PIP_INDEX_URL" --upgrade pip setuptools wheel; \
    "$FORUM_OCR_VENV/bin/pip" install --no-cache-dir \
        --index-url "$PYTORCH_INDEX_URL" \
        torch==2.2.2 \
        torchvision==0.17.2; \
    "$FORUM_OCR_VENV/bin/pip" install --no-cache-dir \
        --index-url "$PIP_INDEX_URL" \
        -r /tmp/ocr-requirements.txt \
        -f "$MODELSCOPE_FIND_LINKS"; \
    if [ "$PRELOAD_MODEL" = "true" ]; then \
        "$FORUM_OCR_VENV/bin/python" /tmp/preload_model.py; \
    fi

WORKDIR /app
COPY --from=builder /out/main /app/main
RUN chmod +x /app/main

CMD ["./main"]
