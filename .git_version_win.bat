@echo off

git rev-parse HEAD > tmp_git_sha1 
set /p GIT_SHA= < tmp_git_sha1 
del tmp_git_sha1 
