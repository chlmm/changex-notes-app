#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 加载 .env 配置文件
if [ -f "$SCRIPT_DIR/.env" ]; then
    # 读取 .env 文件并 export（支持带注释和空行）
    set -a
    while IFS= read -r line || [ -n "$line" ]; do
        # 跳过空行和注释
        [[ -z "$line" || "$line" =~ ^[[:space:]]*# ]] && continue
        # 只处理包含 = 的行
        if [[ "$line" == *"="* ]]; then
            eval "$line"
        fi
    done < "$SCRIPT_DIR/.env"
    set +a
fi

# 处理相对路径（相对于项目根目录）
if [[ "$STATIC_PATH" == ./* ]]; then
    STATIC_PATH="$SCRIPT_DIR/${STATIC_PATH:2}"
fi

# 端口配置
BACKEND_PORT=8080
FRONTEND_PORT=3000

# 查找并杀死占用端口的进程（只杀监听端口的进程）
kill_port() {
    local port=$1
    local name=$2
    
    # 只查找 LISTEN 状态的进程，避免误杀 vscode-server 等
    local pids=($(lsof -t -i:$port -sTCP:LISTEN 2>/dev/null))
    
    if [ ${#pids[@]} -eq 0 ]; then
        return 0
    fi
    
    echo -e "${YELLOW}发现 $name 端口 $port 被占用，正在停止进程...${NC}"
    
    for pid in "${pids[@]}"; do
        # 跳过无效或当前 shell 的 PID
        if [ -z "$pid" ] || [ "$pid" = "$$" ]; then
            continue
        fi
        
        # 先尝试 SIGTERM，再 SIGKILL
        if kill -0 "$pid" 2>/dev/null; then
            kill -15 "$pid" 2>/dev/null
            sleep 0.5
            
            if kill -0 "$pid" 2>/dev/null; then
                kill -9 "$pid" 2>/dev/null
            fi
        fi
    done
    
    sleep 1
    echo -e "${GREEN}已停止端口 $port 的进程${NC}"
}

# 启动后端
start_backend() {
    echo -e "${GREEN}启动后端服务...${NC}"
    cd "$SCRIPT_DIR/notes-server"
    export NOTES_PATH="$NOTES_PATH"
    export STATIC_PATH="$STATIC_PATH"
    go run main.go > /tmp/notes-server.log 2>&1 &
    echo $! > /tmp/notes-server.pid
    sleep 2
    echo -e "${GREEN}后端服务已启动在 http://localhost:$BACKEND_PORT${NC}"
}

# 启动前端
start_frontend() {
    echo -e "${GREEN}启动前端服务...${NC}"
    cd "$SCRIPT_DIR/notes-web"
    npm run dev > /tmp/notes-web.log 2>&1 &
    echo $! > /tmp/notes-web.pid
    sleep 3
    # 检查是否使用了备用端口
    local actual_port=$(grep -oP 'localhost:\K\d+' /tmp/notes-web.log 2>/dev/null | head -1)
    if [ -n "$actual_port" ]; then
        echo -e "${GREEN}前端服务已启动在 http://localhost:$actual_port${NC}"
    else
        echo -e "${GREEN}前端服务已启动${NC}"
    fi
}

# 显示日志
show_logs() {
    echo ""
    echo -e "${YELLOW}===== 后端日志 =====${NC}"
    tail -20 /tmp/notes-server.log 2>/dev/null || echo "暂无日志"
    echo ""
    echo -e "${YELLOW}===== 前端日志 =====${NC}"
    tail -10 /tmp/notes-web.log 2>/dev/null || echo "暂无日志"
}

# 停止所有服务
stop_all() {
    echo -e "${YELLOW}停止所有服务...${NC}"
    kill_port $BACKEND_PORT "后端"
    kill_port $FRONTEND_PORT "前端"
    # 也尝试杀死备用端口
    for port in 3001 3002 3003; do
        kill_port $port "前端备用"
    done
    rm -f /tmp/notes-server.pid /tmp/notes-web.pid 2>/dev/null
    echo -e "${GREEN}所有服务已停止${NC}"
}

# 主函数
main() {
    case "${1:-start}" in
        start)
            echo -e "${GREEN}========== Notes 应用启动 ==========${NC}"
            kill_port $BACKEND_PORT "后端"
            kill_port $FRONTEND_PORT "前端"
            
            start_backend
            start_frontend
            
            echo ""
            echo -e "${GREEN}======================================${NC}"
            echo -e "${GREEN}所有服务已启动！${NC}"
            echo -e "后端: http://localhost:$BACKEND_PORT"
            echo -e "前端: http://localhost:$FRONTEND_PORT (或备用端口)"
            echo -e "${GREEN}======================================${NC}"
            ;;
        stop)
            stop_all
            ;;
        restart)
            stop_all
            sleep 1
            $0 start
            ;;
        logs)
            show_logs
            ;;
        *)
            echo "用法: $0 {start|stop|restart|logs}"
            exit 1
            ;;
    esac
}

main "$@"
