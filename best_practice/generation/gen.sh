#!/bin/bash

set -e

# 四个参数：模板源文件，包名，实际需要具体化的类型，用于构造目标文件名的后缀。
SRC_FILE=${1}
PACKAGE=${2}
TYPE=${3}
DES=${4}
#uppcase the first char
PREFIX="$(tr '[:lower:]' '[:upper:]' <<< ${TYPE:0:1})${TYPE:1}"

DES_FILE=$(echo ${TYPE}| tr '[:upper:]' '[:lower:]')_${DES}.go

# Sed 简明教程：https://coolshell.cn/articles/9104.html
sed 's/PACKAGE_NAME/'"${PACKAGE}"'/g' ${SRC_FILE} | \
    sed 's/GENERIC_TYPE/'"${TYPE}"'/g' | \
    sed 's/GENERIC_NAME/'"${PREFIX}"'/g' > ${DES_FILE}

# 在工程目录中直接执行 go generate 命令，就会生成两份代码：
# - uint32_container.go
# - string_container.go

# 甚至不需要自己手写脚本，直接使用第三方已经写好的工具：
# - Genny     https://github.com/cheekybits/genny
# - Generic   https://github.com/taylorchu/generic
# - GenGen    https://github.com/joeshaw/gengen
# - Gen       https://github.com/clipperhouse/gen