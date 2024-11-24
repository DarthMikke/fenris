millim.devtools.dotenv
=========

Make a .env file from a set of variables.

Role Variables
--------------

**`dest`** - file name of the target file.

**`envvars`** - variables to set in the .env file.

**`mode`** - mode of the target file. Should be a string. Default: `"0664"`.

Example Playbook
----------------

Including an example of how to use your role (for instance, with variables passed in as parameters) is always nice for users too:

    - hosts: servers
      roles:
         - { role: username.rolename, x: 42 }

License
-------

BSD

Author Information
------------------

Michal Jan Warecki

https://www.millim.no/

mwarecki1@gmail.com
