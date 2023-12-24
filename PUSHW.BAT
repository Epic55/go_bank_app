set d="%DATE% - %TIME%" 
git add . 
git commit -m %d% 
git -c http.sslVerify=false push 