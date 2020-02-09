#! /bin/bash

# Check Empty Bucket
curl -ivXGET --user myuser:mypassword localhost:8080

# Upload file
curl -ivXPOST --upload-file word.txt --user myuser:mypassword localhost:8080/word.txt

# Upload file in directory
curl -ivXPOST --upload-file sentence.txt --user myuser:mypassword localhost:8080/testdir/sentence.txt

# Copy file
curl -ivXPUT -d "{\"dist\": \"testdir\/word.txt\", \"recursive\": false}" --user myuser:mypassword localhost:8080/word.txt

# Get directory
curl -ivXGET --user myuser:mypassword localhost:8080/testdir/

# Get copied file
curl -ivXGET --user myuser:mypassword localhost:8080/testdir/word.txt

# Copy directory
curl -ivXPUT -d "{\"dist\": \"test\/\", \"recursive\": true}" --user myuser:mypassword localhost:8080/testdir/

# Get copied directory
curl -ivXGET --user myuser:mypassword localhost:8080/test/

# Get file in copied directory
curl -ivXGET --user myuser:mypassword localhost:8080/test/word.txt

# Delete file
curl -ivXDELETE  --user myuser:mypassword localhost:8080/testdir/word.txt

# Delete directory
curl -ivXDELETE --user myuser:mypassword localhost:8080/testdir/

# Delete copied file
curl -ivXDELETE --user myuser:mypassword localhost:8080/test/word.txt

# Delete copied directory
curl -ivXDELETE --user myuser:mypassword localhost:8080/test/

# Delete original file
curl -ivXDELETE --user myuser:mypassword localhost:8080/word.txt
