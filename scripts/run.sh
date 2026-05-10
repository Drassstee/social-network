#!/bin/bash

# Social Network - Unified Run Script
# This script performs a full clean, build, and start sequence.

# Colors
GREEN="\033[1;32m"
BLUE="\033[1;34m"
ORANGE="\033[1;33m"
RED="\033[1;31m"
NC="\033[0m"

print_status() {
    printf "${GREEN}[+]${NC} ${BLUE}%-15s${NC} ${ORANGE}%s${NC}\n" "$1" "$2"
}

# Ensure we are in the project root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR/.." || exit 1

echo -e "${BLUE}=======================================${NC}"
echo -e "${BLUE}   Social Network Control Center     ${NC}"
echo -e "${BLUE}=======================================${NC}"

# 1. Clean Environment
print_status "Clean" "Purging previous containers and volumes..."
docker compose down --volumes --remove-orphans --rmi local > /dev/null 2>&1

# 2. Build
print_status "Build" "Building Docker images..."
if docker compose build; then
    print_status "Build" "Success."
else
    echo -e "${RED}[!] Build failed.${NC}"
    exit 1
fi

# 3. Start
print_status "Run" "Starting the network..."
if docker compose up -d; then
    print_status "Network" "Operational."
    echo -e "\n${GREEN}Access Points:${NC}"
    echo -e "  - Frontend: ${BLUE}http://localhost${NC}"
    echo -e "  - Backend:  ${BLUE}http://localhost:8080${NC}"
    echo -e "\n${ORANGE}Note: Database has been reset to a fresh state.${NC}\n"
else
    echo -e "${RED}[!] Failed to start containers.${NC}"
    exit 1
fi
