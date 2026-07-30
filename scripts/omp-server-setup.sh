#!/usr/bin/env bash
# 在阿里云服务器上安装、配置 omp，并测量资源消耗
# 默认不走内部梯子；仅在网络失败时尝试 127.0.0.1:7890
set -euo pipefail

PROXY_URL="${OMP_FALLBACK_PROXY:-http://127.0.0.1:7890}"
LOG_DIR="/root/omp-benchmark"
mkdir -p "$LOG_DIR"

log() { echo "[$(date '+%H:%M:%S')] $*"; }

with_proxy() {
  HTTP_PROXY="$PROXY_URL" HTTPS_PROXY="$PROXY_URL" ALL_PROXY="$PROXY_URL" "$@"
}

curl_ok() {
  local url="$1"
  curl -fsSL -o /dev/null --max-time 20 -w "%{http_code}" "$url" 2>/dev/null || echo "000"
}

curl_get() {
  local url="$1" out="$2"
  if curl -fsSL --max-time 30 "$url" -o "$out"; then
    return 0
  fi
  log "直连失败，尝试走梯子: $url"
  if with_proxy curl -fsSL --max-time 30 "$url" -o "$out"; then
    return 0
  fi
  return 1
}

snapshot() {
  local tag="$1"
  local file="$LOG_DIR/${tag}.txt"
  {
    echo "=== $tag @ $(date -Is) ==="
    free -m
    echo "--- top mem ---"
    ps aux --sort=-%mem | head -12
    echo "--- load ---"
    uptime
  } | tee "$file"
}

install_bun() {
  if command -v bun >/dev/null 2>&1; then
    log "bun 已安装: $(bun --version)"
    return 0
  fi

  log "安装 bun..."
  local installer="/tmp/bun-install.sh"
  if ! curl_get "https://bun.sh/install" "$installer"; then
    log "ERROR: 无法下载 bun 安装脚本"
    exit 1
  fi
  bash "$installer"
  export BUN_INSTALL="${BUN_INSTALL:-$HOME/.bun}"
  export PATH="$BUN_INSTALL/bin:$PATH"
  grep -q 'BUN_INSTALL' ~/.bashrc 2>/dev/null || cat >> ~/.bashrc <<'EOF'

# bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"
EOF
  log "bun 安装完成: $(bun --version)"
}

install_omp() {
  export BUN_INSTALL="${BUN_INSTALL:-$HOME/.bun}"
  export PATH="$BUN_INSTALL/bin:$PATH"

  if command -v omp >/dev/null 2>&1; then
    log "omp 已安装: $(omp --version 2>&1 || true)"
    return 0
  fi

  log "尝试官方安装脚本..."
  local installer="/tmp/omp-install.sh"
  if curl_get "https://omp.sh/install" "$installer"; then
    if bash "$installer"; then
      log "官方脚本安装成功"
      return 0
    fi
  fi

  log "官方脚本失败，改用 bun 全局安装..."
  if bun install -g @oh-my-pi/pi-coding-agent; then
    log "bun 全局安装成功"
    return 0
  fi

  log "bun 直连失败，尝试 npm 镜像..."
  if with_proxy bun install -g @oh-my-pi/pi-coding-agent; then
    log "通过梯子安装成功"
    return 0
  fi

  log "ERROR: omp 安装失败"
  exit 1
}

configure_omp() {
  mkdir -p ~/.omp/agent
  local cfg="$HOME/.omp/agent/config.yml"

  if [[ ! -f "$cfg" ]]; then
    cat > "$cfg" <<'EOF'
# 2G 小机器优化配置
tools:
  enabled:
    - github
subagents:
  concurrency: 1
EOF
    log "已写入默认配置: $cfg"
  else
    log "配置已存在，跳过: $cfg"
  fi

  # 不设置 HTTP_PROXY —— 梯子只给其他服务用
  grep -q 'OMP_NO_PROXY_NOTE' ~/.bashrc 2>/dev/null || cat >> ~/.bashrc <<'EOF'

# omp: 默认不走内部梯子（梯子给其他服务用）
# 如需临时启用: export OMP_USE_PROXY=1
OMP_NO_PROXY_NOTE=1
EOF
}

install_gh() {
  if command -v gh >/dev/null 2>&1; then
    log "gh 已安装: $(gh --version | head -1)"
    return 0
  fi

  log "安装 GitHub CLI..."
  if command -v dnf >/dev/null 2>&1; then
  dnf install -y 'dnf-command(config-manager)' 2>/dev/null || true
  dnf config-manager --add-repo https://cli.github.com/packages/rpm/gh-cli.repo 2>/dev/null || true
  if dnf install -y gh 2>/dev/null; then
    log "gh 安装成功"
    return 0
  fi
  fi

  if command -v yum >/dev/null 2>&1; then
    yum install -y gh 2>/dev/null && return 0
  fi

  log "WARN: gh 未安装，自动开 PR 需要后续手动 gh auth login"
}

measure_omp_idle() {
  export BUN_INSTALL="${BUN_INSTALL:-$HOME/.bun}"
  export PATH="$BUN_INSTALL/bin:$PATH"

  log "测量 omp 空闲启动消耗（--version）..."
  local before after
  before=$(free -m | awk '/Mem:/ {print $3}')
  /usr/bin/time -f "TIME: real=%e user=%U sys=%S maxrss=%MKB" \
    omp --version > "$LOG_DIR/omp-version.out" 2> "$LOG_DIR/omp-version.time" || true
  after=$(free -m | awk '/Mem:/ {print $3}')
  echo "mem_delta_version=$((after - before))MB" | tee "$LOG_DIR/omp-version.memdelta"
}

measure_omp_oneshot() {
  export BUN_INSTALL="${BUN_INSTALL:-$HOME/.bun}"
  export PATH="$BUN_INSTALL/bin:$PATH"

  if [[ -z "${ANTHROPIC_API_KEY:-}" && -z "${OPENAI_API_KEY:-}" && -z "${OPENROUTER_API_KEY:-}" ]]; then
    log "WARN: 未检测到 API Key，跳过 omp -p 实测（仅完成安装与版本检测）"
  log "如需实测推理消耗，请设置: export ANTHROPIC_API_KEY=sk-... 后重跑 measure 段"
    return 0
  fi

  log "运行 omp -p 轻量任务（list files）..."
  local before after
  before=$(free -m | awk '/Mem:/ {print $3}')
  /usr/bin/time -f "TIME: real=%e user=%U sys=%S maxrss=%MKB" \
    omp -p "只列出当前目录下的文件，不要修改任何东西，完成后直接退出" \
    > "$LOG_DIR/omp-oneshot.out" 2> "$LOG_DIR/omp-oneshot.time" || true
  after=$(free -m | awk '/Mem:/ {print $3}')
  echo "mem_delta_oneshot=$((after - before))MB" | tee "$LOG_DIR/omp-oneshot.memdelta"

  log "查找 omp 相关进程..."
  ps aux | grep -E '[o]mp|[b]un.*pi-coding' | tee "$LOG_DIR/omp-processes.txt" || true
}

main() {
  log "开始 omp 安装与压测（默认不走梯子）"
  log "网络探测:"
  log "  omp.sh -> $(curl_ok https://omp.sh/)"
  log "  github.com -> $(curl_ok https://github.com)"
  log "  registry.npmjs.org -> $(curl_ok https://registry.npmjs.org)"

  snapshot "before-install"
  install_bun
  install_omp
  configure_omp
  install_gh
  snapshot "after-install"

  export BUN_INSTALL="${BUN_INSTALL:-$HOME/.bun}"
  export PATH="$BUN_INSTALL/bin:$PATH"
  log "omp 版本: $(omp --version 2>&1 || echo unknown)"

  measure_omp_idle
  measure_omp_oneshot
  snapshot "after-benchmark"

  log "完成。日志目录: $LOG_DIR"
  ls -la "$LOG_DIR"
}

main "$@"
