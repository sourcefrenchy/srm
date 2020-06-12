# Home-made secure rm for OSX
A secure delete for OSX (and other platforms if you want to cross-compile, courtesy of using Golang).
Reading on secure delete techniques and decided to set on NSA 130-2, i.e. 3 passes: random, random, verify last random.

# How?
Well, I just use rand.Seed(time.Now().UTC().UnixNano()) and rand.Read() for each random pass.
Then I open the file before doing os.delete on it to ensure the last random was written correctly via [Blake](https://en.wikipedia.org/wiki/BLAKE_(hash_function)) 256-bits.
You can also use "rm -P" from OSX since srm is not longer available:
```     
     -P          Overwrite regular files before deleting them.  Files are overwrit-
                 ten three times, first with the byte pattern 0xff, then 0x00, and
                 then 0xff again, before they are deleted.
```
                 
# Why?
Well I didnt know the wheel existed, wrote my mini srm before deciding to read the man rm... But I am still trying to learn Golang so still fun to do.

# Feedbacks
Yes, I am really really learning Golang coming from Python so please PR/help/feedback really welcome on making better code, more readable, more efficient and secure if need be :)
