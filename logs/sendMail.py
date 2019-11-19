# -*- coding:utf-8 -*-
'''
  update: 2019-10-23
  function: send feedback and log to mail using python
'''

import yagmail

# 附件
prefix = "/Pages/logs/"
file = prefix + "feedback.txt"
logfile = prefix + "Pages.log"


#链接邮箱服务器  
yag = yagmail.SMTP(user="", password="", host='')

# 邮箱正文
contents = "请查收查老师系统日志及近期反馈附件"
title = "Pages日志及反馈"

# send
yag.send('xxxxx@xx.xx', title, contents, [file, logfile])

print("send successfully.")

