# tiger
**Go mini-web framework**

**How to use tiger framework**

*Step 1*
Ensure config/config.json are setup correctly for your environment.

*Step 2*
Start to add your application specific code in util/http/handler_util.go Refer to the relevant package documentation on how to do it.

*Step 3*
Compile by running go build. tiger.exe or tiger will be created.

*Step 4*
From Windows Command Prompt or Linux terminal, execute tiger.exe or tiger &

*Step 5*
Use a browser and navigate to your configured url in Step 1 config.json e.g http://localhost:8000
You should see a message I am alive! This mean your http server is up and running.
To shutdown, send a SIGINT signal. Ctrl-C for Windows Command Prompt. kill -SIGINT <pid> for Linux.
  
**Install below Go dependency packages separately**

*Step 1*
go get -v github.com/go-sql-driver/mysql

*Step 2*
go get -v golang.org/x/text

**View the framework documentation**

*Step 1*
godoc -http=localhost:6060

*Step 2*
Use a browser and navigate to http://localhost:6060/pkg/tiger/

**Contact**

Any bug/suggestion/feedback can mail to sohguanh@gmail.com

**Donate**

Yes I know almost 99.99% of people do not donate. But I still include a PayPal Donate for that 0.01% (if they exist) and they appreciate the work I have done or my work has helped them in anyway.

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://paypal.me/sohguanh/1)

:tiger: :tiger: :tiger:
