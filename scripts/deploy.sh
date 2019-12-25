#make all
go build -o nestingbot
upx nestingbot
rsync -avzrP --exclude 'data/*' ./ simplecloud:/home/god/sites/nestingbot/
ssh simplecloud "sudo stop nestingbot && sudo start nestingbot"
