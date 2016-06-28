# dplaravel

**dplaravel** creates a build folder in your project, run any command you issue and then creates a tar.gz file.

> I developed DPLaravel in some hours at night, so don't expect nothing more than a super-simple tool (But maybe we can move it towards!)


### Usage

it cannot be simpler than download **dplaravel** and add a **dplaravel.json** file on project root, then run it on bash:

```
dplaravel
```

The dplaravel.json should be something like this:

```
{
  "name": "project-name",
  "branch": "master",
  "dest": "build",
  "commands": [
    { "cmd": "composer install --no-ansi --no-dev --no-interaction --no-scripts --optimize-autoloader" },
    { "cmd": "gulp --production" }
  ],
  "sftp": {
    "host": "",
    "user": "user",
    "password": "password",
    "directory": "/var/www"
  }
}
```

**Wow**

**dplaravel** is build on Go, so take a look at main.go. You can build it on your own too.


#### TODO

- Deploy fo SFTP
