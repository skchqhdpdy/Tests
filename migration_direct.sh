ARCHIVE_PATH="/home/$(whoami)/migration_direct.tar.gz"
read -p ".git-credentials .gitconfig bancho.py/ certs/ guweb/ tools/ venv/ 해당 파일 및 폴더를 ${ARCHIVE_PATH} 로 압축 하시겠습니까? (y/n, (y)) : " isArchive
if [[ $isArchive == "n" ]]; then
    read -p "${ARCHIVE_PATH} 를 압축 해제 하시겠습니까? (y/n, (y)) : " isDearchive
    if [[ $isDearchive == "n" ]]; then
        echo "압축 해제 취소"; exit
    fi
    cmd="sudo tar xvpzf \"${ARCHIVE_PATH}\""
    echo "${cmd}"; eval $cmd
    echo "exit"; exit
fi
if [ -f "${ARCHIVE_PATH}" ]; then
    read -p "${ARCHIVE_PATH} 파일이 이미 존재합니다. 덮어쓰시겠습니까? (y/n, (y)) : " isOverwrite
    if [[ $isOverwrite == "n" ]]; then
        echo "압축 취소"; exit
    fi
fi
cmd="sudo tar cvpzf \"${ARCHIVE_PATH}\" .git-credentials .gitconfig bancho.py/ certs/ guweb/ tools/ venv/"
echo "${cmd}"; eval $cmd
echo "exit"; exit