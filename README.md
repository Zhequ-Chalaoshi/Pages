# Pages
**P**ublic-shared     
**A**nonymous      
**G**eneral-teachers'      
**E**valuation      
**S**ystem     



## 介绍
支持部分本科生教师和研究生导师的查询与评论。  
本科生初始数据来源于查老师，研究生导师数据来源于导师评价网及个人主页。  
数据库已上传到[此处](https://box.zjuqsc.com/-41188957)，可自行下载使用。




## 技术栈
* 前端：html + bootstrap + js  
* 后端：[Golang](https://golang.org) + [beego框架](https://beego.me)
* 数据库：MySQL




## 运行环境
* Go编译器（1.12版本及以上）
* Python3解释器
### 1. 安装依赖包
```shell
# for beego
go get github.com/astaxie/beego
go get github.com/beego/bee
# for python 
conda install yagmail
## or
pip install yagmail
```

### 2. 运行
```shell
cd $Pages
bee run
```
### 3. 部署
```shell
# 打包成Linux可执行文件
bee pack -be GOOS=linux

# 向系统添加定期执行任务(定期发送log和feedback邮件)
crontab -e 
22 22 * * 2 /bin/python $Pages/logs/sendMail.py
```



## TODO

* 点赞功能
* 按课程查询
* 自动添加教师
* 界面美化与设计
* 前后端分离实现
* 智能化垃圾评论筛选
* 网站安全性
* ...



## Contribution

Fork this repository, write your own code and submit a `pull request`.



## Reference
[导师评价网](https://mysupervisor.org)  
[研究生老师](http://yjsds.zju.edu.cn/daoshiInfoList.jsp)   
[浙大个人主页](https://person.zju.edu.cn/index/)   
[前人所做的查老师](https://github.com/ssaichixbg/chalaoshi)
