# To Do List and Notes

- Need to add a way to search for packages.
- use PA 2757  and 491 to test "dual packaging"
- keeup up to 3-5 last versions.  Have a pointer for latest and current.
  - Allow user to control what current points to.
- Going to use <https://github.com/gobuffalo/buffalo> for the web bits.
- Need a package to work with the app package
  - unpack package
  - create a copy for each type: idx, shc, fwd - as needed
  - parse the directory structure into something
  - identify specific files that need to be removed or edited for this server type
  - take appropriate actions
  - repackage files into their own archive for deployment
- Put some logic into the basic package download
  - Check for redirect.
  - If it redirects to github, maybe we parse out the DL link from that page?
- Add switch/logic to deal with github links?
- Maybe switch for direct tgz links ?
