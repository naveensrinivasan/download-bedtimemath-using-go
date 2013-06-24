## Download bedtime math using go ##

[http://bedtimemath.org/](http://bedtimemath.org/ "Bedtime math")

My kids love to read bedtime math and would prefer reading from the actual 
printed books more than on a screen.

Bedtime math has a great archive of math stories and riddles. 

This was a great opportunity to try out something for my own family using GO. 

The code download's the pages 1-100 (configurable) and looks for the 
actual post url within the page and then request the actual post as pdf using 
the convertapi and then saves it.  

It depends on **github.com/moovweb/gokogiri** for xpath.

**ConvertAPI** for converting webpage to PDF http://do.convertapi.com/Web2Pdf 

All of these happens concurrently using ***channels***.    

Here is the first [100](https://github.com/naveensrinivasan/download-bedtimemath-using-go/blob/master/bedtimemath.pdf "First 100 posts ") posts as pdf. 
 


