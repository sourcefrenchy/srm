# srm
A secure delete for OSX (and other if you want to cross-compile).
Reading on techniques and set on NSA 130-2, i.e. 3 passes: random, random, verify last random

You can also use "rm -P" from OSX since srm is not longer available:
```     
     -P          Overwrite regular files before deleting them.  Files are overwrit-
                 ten three times, first with the byte pattern 0xff, then 0x00, and
                 then 0xff again, before they are deleted.
```
                 
# Why reinventing the wheel?
Well I didnt know the wheel existed, wrote my mini srm before deciding to read the man rm... But I am still trying to learn Golang so still fun to do.
