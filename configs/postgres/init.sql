CREATE DATABASE testchat;
CREATE user testchat with password 'testchat';
GRANT ALL PRIVILEGES ON DATABASE testchat TO testchat;

CREATE DATABASE chatapp;
CREATE user chatapp with password 'chatapp';
GRANT ALL PRIVILEGES ON DATABASE chatapp TO chatapp;