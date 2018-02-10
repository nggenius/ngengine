echo "begin reverse"
del .\model\* /s/q
xorm_tool reverse -s mysql root:@tcp(127.0.0.1:3306)/nx_base?charset=utf8 ./templates/gongengine
pause