#!/bin/bash

set -eu
# don't set any options, otherwise it will fail arbitrarily
# set -euo pipefail

work_dir=$1
repo_name=$2
pr_num=$3
origin_repo_url=$4
remote_repo_url=$5

set +e
test -d $work_dir || mkdir -p $work_dir
set -e

cd $work_dir

if [ ! -d $repo_name ]; then
    git clone -q $origin_repo_url
fi

cd $repo_name

git fetch origin "pull/$pr_num/head"

git checkout "FETCH_HEAD"

branch="pull$pr_num"

git checkout -b $branch

git push -f $remote_repo_url $branch
